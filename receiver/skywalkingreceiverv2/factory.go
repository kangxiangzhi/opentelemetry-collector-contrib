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
	"net"
	"strconv"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
)

const (
	typeStr                 = "skywalkingv2"
	protoGRPCV3             = "grpc"
	defaultGRPCBindEndpoint = "0.0.0.0:11800"
)

var receivers = map[*Config]*SkyWalkingReceiverPhp{}

//NewFactory creates a new skyWalking receiver factory
func NewFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesReceiver(createTraceReceiver),
		component.WithMetricsReceiver(createMetricsReceiver))
}

// CreateDefaultConfig creates the default configuration for Jaeger receiver.
func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
		Protocols: Protocols{
			GrpcV3: &configgrpc.GRPCServerSettings{
				NetAddr: confignet.NetAddr{
					Endpoint:  defaultGRPCBindEndpoint,
					Transport: "tcp",
				},
			},
		},
	}
}

//createTraceReceiver creates a trace receiver based on provided config.
func createTraceReceiver(
	ctx context.Context,
	set component.ReceiverCreateSettings,
	cfg config.Receiver,
	nextConsumer consumer.Traces,
) (component.TracesReceiver, error) {
	return createReceiver(cfg, nextConsumer, nil, set)
}

func createMetricsReceiver(
	ctx context.Context,
	set component.ReceiverCreateSettings,
	cfg config.Receiver,
	nextConsumer consumer.Metrics,
) (component.MetricsReceiver, error) {
	return createReceiver(cfg, nil, nextConsumer, set)
}

func createReceiver(cfg config.Receiver,
	nextTraceConsumer consumer.Traces,
	nextMetricConsumer consumer.Metrics,
	params component.ReceiverCreateSettings) (*SkyWalkingReceiverPhp, error) {
	rCfg := cfg.(*Config)
	receiver, ok := receivers[rCfg]
	config := resetConfiguration(rCfg)
	if !ok {
		var err error
		// We don't have a receiver, so create one.
		receiver, err = newSkyWalkingReceiver(rCfg.ID(), nextTraceConsumer, nextMetricConsumer, config, params)
		if err != nil {
			return nil, err
		}
		// Remember the receiver in the map
		receivers[rCfg] = receiver
	} else {
		if nextTraceConsumer != nil {
			receiver.nextTraceConsumer = nextTraceConsumer
		}
		if nextMetricConsumer != nil {
			receiver.nextMetricConsumer = nextMetricConsumer
		}
	}
	return receiver, nil
}

func resetConfiguration(rCfg *Config) *configuration {
	var config configuration
	if rCfg.Protocols.GrpcV3 != nil {
		var err error
		config.CollectorGRPCV3Port, err = extractPortFromEndpoint(rCfg.Protocols.GrpcV3.NetAddr.Endpoint)
		if err != nil {
			return nil
		}
		config.CollectorGRPCServerSettings = *rCfg.Protocols.GrpcV3
		config.ComponentMap = extractComponentMapFromComponents(rCfg.Components.Components)
	}
	return &config
}

// extract the port number from string in "address:port" format. If the
// port number cannot be extracted returns an error.
func extractPortFromEndpoint(endpoint string) (int, error) {
	_, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return 0, fmt.Errorf("endpoint is not formatted correctly: %s", err.Error())
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("endpoint port is not a number: %s", err.Error())
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port number must be between 1 and 65535")
	}
	return int(port), nil
}

func extractComponentMapFromComponents(components []*Component) map[int32]string {
	componentMap := make(map[int32]string, len(components))
	for _, componentInfo := range components {
		componentMap[componentInfo.ComponentId] = componentInfo.ComponentName
	}
	return componentMap
}
