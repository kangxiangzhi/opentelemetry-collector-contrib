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
	"fmt"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
)

const (
	protocolsFieldName = "protocols"
)

//Protocols
type Protocols struct {
	GrpcV3 *configgrpc.GRPCServerSettings `mapstructure:"grpc"`
}

//components
type Component struct {
	ComponentName string `mapstructure:"name"`
	ComponentId   int32  `mapstructure:"id"`
}

type Components struct {
	Components []*Component `mapstructure:"components"`
}

//Config defines configuration for SkyWalking receiver.
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"`
	Protocols               `mapstructure:"protocols"`
	Components              `mapstructure:",squash"`
}

var _ config.Receiver = (*Config)(nil)
var _ config.Unmarshallable = (*Config)(nil)

func (cfg *Config) Validate() error {

	if cfg.GrpcV3 == nil {
		return fmt.Errorf("must specify at least one protocol when using the skywalking receiver")
	}

	return nil
}

func (cfg *Config) Unmarshal(componentParser *config.Map) error {
	if componentParser == nil || len(componentParser.AllKeys()) == 0 {
		return fmt.Errorf("empty config for Skywalking receiver")
	}

	// UnmarshalExact will not set struct properties to nil even if no key is provided,
	// so set the protocol structs to nil where the keys were omitted.
	err := componentParser.UnmarshalExact(cfg)
	if err != nil {
		return err
	}

	protocols, err := componentParser.Sub(protocolsFieldName)
	if err != nil {
		return err
	}

	knownProtocols := 0
	if !protocols.IsSet(protoGRPCV3) {
		cfg.GrpcV3 = nil
	} else {
		knownProtocols++
	}
	if len(protocols.AllKeys()) != knownProtocols {
		return fmt.Errorf("unknown protocols in the skywalking receiver")
	}

	return nil
}
