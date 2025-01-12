// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v3

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	v3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TraceSegmentReportServiceClient is the client API for TraceSegmentReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TraceSegmentReportServiceClient interface {
	Collect(ctx context.Context, opts ...grpc.CallOption) (TraceSegmentReportService_CollectClient, error)
}

type traceSegmentReportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTraceSegmentReportServiceClient(cc grpc.ClientConnInterface) TraceSegmentReportServiceClient {
	return &traceSegmentReportServiceClient{cc}
}

func (c *traceSegmentReportServiceClient) Collect(ctx context.Context, opts ...grpc.CallOption) (TraceSegmentReportService_CollectClient, error) {
	stream, err := c.cc.NewStream(ctx, &TraceSegmentReportService_ServiceDesc.Streams[0], "/TraceSegmentReportService/collect", opts...)
	if err != nil {
		return nil, err
	}
	x := &traceSegmentReportServiceCollectClient{stream}
	return x, nil
}

type TraceSegmentReportService_CollectClient interface {
	Send(*agentv3.SegmentObject) error
	CloseAndRecv() (*v3.Commands, error)
	grpc.ClientStream
}

type traceSegmentReportServiceCollectClient struct {
	grpc.ClientStream
}

func (x *traceSegmentReportServiceCollectClient) Send(m *agentv3.SegmentObject) error {
	return x.ClientStream.SendMsg(m)
}

func (x *traceSegmentReportServiceCollectClient) CloseAndRecv() (*v3.Commands, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(v3.Commands)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TraceSegmentReportServiceServer is the server API for TraceSegmentReportService service.
// All implementations must embed UnimplementedTraceSegmentReportServiceServer
// for forward compatibility
type TraceSegmentReportServiceServer interface {
	Collect(TraceSegmentReportService_CollectServer) error
	mustEmbedUnimplementedTraceSegmentReportServiceServer()
}

// UnimplementedTraceSegmentReportServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTraceSegmentReportServiceServer struct {
}

func (UnimplementedTraceSegmentReportServiceServer) Collect(TraceSegmentReportService_CollectServer) error {
	return status.Errorf(codes.Unimplemented, "method Collect not implemented")
}
func (UnimplementedTraceSegmentReportServiceServer) mustEmbedUnimplementedTraceSegmentReportServiceServer() {
}

// UnsafeTraceSegmentReportServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TraceSegmentReportServiceServer will
// result in compilation errors.
type UnsafeTraceSegmentReportServiceServer interface {
	mustEmbedUnimplementedTraceSegmentReportServiceServer()
}

func RegisterTraceSegmentReportServiceServer(s grpc.ServiceRegistrar, srv TraceSegmentReportServiceServer) {
	s.RegisterService(&TraceSegmentReportService_ServiceDesc, srv)
}

func _TraceSegmentReportService_Collect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TraceSegmentReportServiceServer).Collect(&traceSegmentReportServiceCollectServer{stream})
}

type TraceSegmentReportService_CollectServer interface {
	SendAndClose(*v3.Commands) error
	Recv() (*agentv3.SegmentObject, error)
	grpc.ServerStream
}

type traceSegmentReportServiceCollectServer struct {
	grpc.ServerStream
}

func (x *traceSegmentReportServiceCollectServer) SendAndClose(m *v3.Commands) error {
	return x.ServerStream.SendMsg(m)
}

func (x *traceSegmentReportServiceCollectServer) Recv() (*agentv3.SegmentObject, error) {
	m := new(agentv3.SegmentObject)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TraceSegmentReportService_ServiceDesc is the grpc.ServiceDesc for TraceSegmentReportService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TraceSegmentReportService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TraceSegmentReportService",
	HandlerType: (*TraceSegmentReportServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "collect",
			Handler:       _TraceSegmentReportService_Collect_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "Tracing.proto",
}
