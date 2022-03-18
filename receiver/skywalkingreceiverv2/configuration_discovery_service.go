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

	configurationV3 "skywalking.apache.org/repo/goapi/collect/agent/configuration/v3"
	v3 "skywalking.apache.org/repo/goapi/collect/common/v3"
)

//Provide query agent dynamic configuration, through the gRPC protocol,
type ConfigurationDiscoveryService struct {
	configurationV3.UnimplementedConfigurationDiscoveryServiceServer
}

//Process the request for querying the dynamic configuration of the agent.
//If there is agent dynamic configuration information corresponding to the service,
//the ConfigurationDiscoveryCommand is returned to represent the dynamic configuration information.
func (cds *ConfigurationDiscoveryService) FetchConfigurations(context.Context,
	*configurationV3.ConfigurationSyncRequest) (*v3.Commands, error) {
	return &v3.Commands{}, nil
}
