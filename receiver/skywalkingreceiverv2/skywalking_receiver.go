// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package skywalkingreceiverv2

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"net"
	"sync"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalkingv2"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiverv2/metrics"
	agentV3Php "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiverv2/php/collect/language/agent/v3"
	managementV3Php "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiverv2/php/collect/management/v3"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiverv2/trace"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	configurationV3 "skywalking.apache.org/repo/goapi/collect/agent/configuration/v3"
	agentV3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	managementV3 "skywalking.apache.org/repo/goapi/collect/management/v3"
)

type configuration struct {
	CollectorGRPCV3Port         int
	ComponentMap                map[int32]string
	CollectorGRPCServerSettings configgrpc.GRPCServerSettings
}

// Receiver type is used to receive spans that were originally intended to be sent to skyWalking.
// This receiver is basically a skyWalking collector.
type SkyWalkingReceiver struct {
	mu sync.Mutex

	nextTraceConsumer  consumer.Traces
	nextMetricConsumer consumer.Metrics
	id                 config.ComponentID
	config             *configuration
	settings           component.ReceiverCreateSettings

	grpcv3 *grpc.Server

	traceReceiver   *trace.ReceiverPhp
	metricsReceiver *metrics.ReceiverPhp

	startOnce sync.Once
	stopOnce  sync.Once

	logger *zap.Logger
}

type SkyWalkingReceiverPhp struct {
	*SkyWalkingReceiver
}

//newSkyWalkingReceiver creates a TracesReceiver that receives traffic as a skyWalking collector
func newSkyWalkingReceiver(
	id config.ComponentID,
	nextTraceConsumer consumer.Traces,
	nextMetricConsumer consumer.Metrics,
	config *configuration,
	set component.ReceiverCreateSettings,
) (*SkyWalkingReceiverPhp, error) {
	r := &SkyWalkingReceiverPhp{
		&SkyWalkingReceiver{
			id:                 id,
			nextTraceConsumer:  nextTraceConsumer,
			nextMetricConsumer: nextMetricConsumer,
			config:             config,
			settings:           set,
		},
	}
	return r, nil
}

// Start runs the trace receiver on the gRPC server. Currently
func (swr *SkyWalkingReceiver) Start(_ context.Context, host component.Host) error {
	swr.mu.Lock()
	defer swr.mu.Unlock()

	var err error

	swr.startOnce.Do(func() {
		err = swr.startGRPCServer(host)
		err = nil
	})
	return err
}

// Shutdown is a method to turn off receiving.
func (swr *SkyWalkingReceiver) Shutdown(context.Context) error {
	var err error
	swr.stopOnce.Do(func() {
		swr.mu.Lock()
		defer swr.mu.Unlock()
		if swr.grpcv3 != nil {
			swr.grpcv3.GracefulStop()
		}
		err = nil
	})
	return err
}

func (swr *SkyWalkingReceiver) startGRPCServer(host component.Host) error {
	if swr.collectorGRPCV3Enabled() {

		opts, err := swr.config.CollectorGRPCServerSettings.ToServerOption(host, swr.settings.TelemetrySettings)
		if err != nil {
			return fmt.Errorf("failed to build the options for the Skywalking gRPC Collector: %v", err)
		}
		swr.grpcv3 = grpc.NewServer(opts...)

		gAddr := swr.collectorGRPCV3Addr()
		gln, gErr := net.Listen("tcp", gAddr)
		if gErr != nil {
			return fmt.Errorf("failed to bind to gRPC address %q: %v", gAddr, gErr)
		}
		managementV3.RegisterManagementServiceServer(swr.grpcv3, &ManagementServiceServer{})
		managementV3Php.RegisterManagementServiceServer(swr.grpcv3, &ManagementServiceServerPhp{})

		configurationV3.RegisterConfigurationDiscoveryServiceServer(swr.grpcv3, &ConfigurationDiscoveryService{})

		skywalkingv2.ComponentMap = swr.config.ComponentMap

		//注册Rpc
		if swr.nextMetricConsumer != nil {
			if err = swr.registerMetricsConsumer(swr.nextMetricConsumer, swr.id); err != nil {
				return err
			}
		}
		if swr.nextTraceConsumer != nil {
			if err = swr.registerTraceConsumer(swr.nextTraceConsumer, swr.id); err != nil {
				return err
			}
		}

		go func() {
			if err := swr.grpcv3.Serve(gln); err != nil {
				host.ReportFatalError(err)
			}
		}()
	}
	return nil
}

func (swr *SkyWalkingReceiver) registerTraceConsumer(tc consumer.Traces,
	id config.ComponentID) error {
	if tc == nil {
		return componenterror.ErrNilNextConsumer
	}
	swr.traceReceiver = trace.New(id, tc, swr.settings, swr.logger)
	if swr.grpcv3 != nil {
		agentV3.RegisterTraceSegmentReportServiceServer(swr.grpcv3, swr.traceReceiver.Receiver)
		agentV3Php.RegisterTraceSegmentReportServiceServer(swr.grpcv3, swr.traceReceiver)
	}
	return nil
}

func (swr *SkyWalkingReceiver) registerMetricsConsumer(mc consumer.Metrics,
	id config.ComponentID) error {
	if mc == nil {
		return componenterror.ErrNilNextConsumer
	}
	swr.metricsReceiver = metrics.New(id, mc, swr.settings, swr.logger)
	if swr.grpcv3 != nil {
		agentV3.RegisterJVMMetricReportServiceServer(swr.grpcv3, swr.metricsReceiver.Receiver)
		agentV3Php.RegisterJVMMetricReportServiceServer(swr.grpcv3, swr.metricsReceiver)
	}
	return nil
}

func (swr *SkyWalkingReceiver) collectorGRPCV3Addr() string {
	var port int
	if swr.config != nil {
		port = swr.config.CollectorGRPCV3Port
	}
	return fmt.Sprintf(":%d", port)
}

func (swr *SkyWalkingReceiver) collectorGRPCV3Enabled() bool {
	return swr.config != nil && swr.config.CollectorGRPCV3Port > 0
}
