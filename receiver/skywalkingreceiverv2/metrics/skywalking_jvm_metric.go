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

package metrics

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalkingv2"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/obsreport"

	agentV3Php "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiverv2/php/collect/language/agent/v3"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	commonV3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentV3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const (
	AuthenticationHeader = "authentication"
	grpcTransport        = "grpc"
	protobufFormat       = "protobuf"
)

//Define the JVM metrics report service
type Receiver struct {
	id                      config.ComponentID
	nextConsumer            consumer.Metrics
	logger                  *zap.Logger
	skywalkingToOtlpMetrics *skywalkingv2.SkywalkingToOtlpMetrics
	agentV3.UnimplementedJVMMetricReportServiceServer
	grpcObsrecv *obsreport.Receiver
}

type ReceiverPhp struct {
	*Receiver
	agentV3Php.UnimplementedJVMMetricReportServiceServer
}

//New
func New(id config.ComponentID, nextConsumer consumer.Metrics, set component.ReceiverCreateSettings, logger *zap.Logger) *ReceiverPhp {
	r := &Receiver{
		id:                      id,
		nextConsumer:            nextConsumer,
		logger:                  logger,
		skywalkingToOtlpMetrics: &skywalkingv2.SkywalkingToOtlpMetrics{Logger: logger},
		grpcObsrecv: obsreport.NewReceiver(obsreport.ReceiverSettings{
			ReceiverID:             id,
			Transport:              grpcTransport,
			ReceiverCreateSettings: set,
		}),
	}
	return &ReceiverPhp{
		Receiver: r,
	}
}

//Collect
func (sjr *Receiver) Collect(ctx context.Context,
	jvmMetricCollector *agentV3.JVMMetricCollection) (*commonV3.Commands, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if ok {
		authentication := headers.Get(AuthenticationHeader)
		if len(authentication) > 0 {
			ctx = sjr.grpcObsrecv.StartTracesOp(ctx)
			pm := sjr.skywalkingToOtlpMetrics.JVMToMetrics(authentication[0], jvmMetricCollector)
			err := sjr.nextConsumer.ConsumeMetrics(ctx, pm)
			sjr.grpcObsrecv.EndTracesOp(ctx, protobufFormat, len(jvmMetricCollector.Metrics), err)
		}
	}
	return &commonV3.Commands{}, nil
}

func (sjr *ReceiverPhp) Collect(ctx context.Context, jvmMetricCollector *agentV3.JVMMetricCollection) (*commonV3.Commands, error) {
	return sjr.Receiver.Collect(ctx, jvmMetricCollector)
}
