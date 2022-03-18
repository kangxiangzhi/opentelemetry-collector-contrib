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

package trace

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/skywalkingv2"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/obsreport"
	"io"

	agentV3Php "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiverv2/php/collect/language/agent/v3"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
	commonV3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentV3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

type Receiver struct {
	id           config.ComponentID
	nextConsumer consumer.Traces
	logger       *zap.Logger
	agentV3.UnimplementedTraceSegmentReportServiceServer

	grpcObsrecv *obsreport.Receiver
}

type ReceiverPhp struct {
	*Receiver
	agentV3Php.UnimplementedTraceSegmentReportServiceServer
}

const (
	grpcTransport  = "grpc"
	protobufFormat = "protobuf"
	authority      = ":authority"
	serviceName    = "serviceName"
)

// New creates a new Receiver reference.
func New(id config.ComponentID,
	nextConsumer consumer.Traces,
	set component.ReceiverCreateSettings,
	logger *zap.Logger) *ReceiverPhp {
	r := &ReceiverPhp{
		Receiver: &Receiver{
			id:           id,
			nextConsumer: nextConsumer,
			logger:       logger,
			grpcObsrecv: obsreport.NewReceiver(obsreport.ReceiverSettings{
				ReceiverID:             id,
				Transport:              grpcTransport,
				ReceiverCreateSettings: set,
			}),
		},
	}
	return r
}

//Collect streaming for client. skyWalking v3 interface
func (swr *Receiver) Collect(tc agentV3.TraceSegmentReportService_CollectServer) error {
	for {
		segmentObject, err := tc.Recv() //segmentObject
		if err == io.EOF {
			tc.SendAndClose(&commonV3.Commands{})
			return nil
		}
		if err != nil {
			swr.logger.Error("Error recv grpc", zap.Error(err))
			return err
		}
		ctx := tc.Context()
		ctx = swr.grpcObsrecv.StartTracesOp(ctx)
		td := skywalkingv2.SegmentToInternalTraces(ctx, segmentObject)
		err = swr.nextConsumer.ConsumeTraces(ctx, td)
		swr.grpcObsrecv.EndTracesOp(ctx, protobufFormat, len(segmentObject.Spans), err)
	}
	return nil
}

//CollectInSync An alternative for trace report by using gRPC unary.skyWalking v3 interface
func (swr *Receiver) CollectInSync(ctx context.Context, segments *agentV3.SegmentCollection) (
	*commonV3.Commands, error) {
	ctx = swr.grpcObsrecv.StartTracesOp(ctx)
	td := skywalkingv2.SegmentCollectionToInternalTraces(ctx, segments)
	err := swr.nextConsumer.ConsumeTraces(ctx, td)
	swr.grpcObsrecv.EndTracesOp(ctx, protobufFormat, len(segments.Segments), err)
	return &commonV3.Commands{}, nil
}

func (swr *ReceiverPhp) Collect(tc agentV3Php.TraceSegmentReportService_CollectServer) error {
	return swr.Receiver.Collect(tc)
}
