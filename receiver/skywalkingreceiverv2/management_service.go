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

	agentV3Php "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiverv2/php/collect/management/v3"
	commonV3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentV3 "skywalking.apache.org/repo/goapi/collect/management/v3"
)

//Define the service reporting the extra information of the instance
type ManagementServiceServer struct {
	agentV3.UnimplementedManagementServiceServer
}

type ManagementServiceServerPhp struct {
	*ManagementServiceServer
	agentV3Php.UnimplementedManagementServiceServer
}

//Report custom properties of a service instance.
func (mss *ManagementServiceServer) ReportInstanceProperties(context.Context,
	*agentV3.InstanceProperties) (*commonV3.Commands, error) {
	return &commonV3.Commands{}, nil
}

//Keep the instance alive in the backend analysis.
//Only recommend to do separate keepAlive report when no trace and metrics needs to be reported.
// Otherwise, it is duplicated.
func (mss *ManagementServiceServer) KeepAlive(context.Context,
	*agentV3.InstancePingPkg) (*commonV3.Commands, error) {
	return &commonV3.Commands{}, nil
}

//Report custom properties of a service instance.
func (mss *ManagementServiceServerPhp) ReportInstanceProperties(context.Context,
	*agentV3.InstanceProperties) (*commonV3.Commands, error) {
	return &commonV3.Commands{}, nil
}

//Keep the instance alive in the backend analysis.
//Only recommend to do separate keepAlive report when no trace and metrics needs to be reported.
// Otherwise, it is duplicated.
func (mss *ManagementServiceServerPhp) KeepAlive(context.Context,
	*agentV3.InstancePingPkg) (*commonV3.Commands, error) {
	return &commonV3.Commands{}, nil
}
