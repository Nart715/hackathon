// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.6.1
// source: proto/betting/betting.proto

package betting

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BettingService_CreateBetting_FullMethodName = "/betting.BettingService/CreateBetting"
)

// BettingServiceClient is the client API for BettingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BettingServiceClient interface {
	CreateBetting(ctx context.Context, in *BettingMessageRequest, opts ...grpc.CallOption) (*BettingResponse, error)
}

type bettingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBettingServiceClient(cc grpc.ClientConnInterface) BettingServiceClient {
	return &bettingServiceClient{cc}
}

func (c *bettingServiceClient) CreateBetting(ctx context.Context, in *BettingMessageRequest, opts ...grpc.CallOption) (*BettingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BettingResponse)
	err := c.cc.Invoke(ctx, BettingService_CreateBetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BettingServiceServer is the server API for BettingService service.
// All implementations must embed UnimplementedBettingServiceServer
// for forward compatibility.
type BettingServiceServer interface {
	CreateBetting(context.Context, *BettingMessageRequest) (*BettingResponse, error)
	mustEmbedUnimplementedBettingServiceServer()
}

// UnimplementedBettingServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBettingServiceServer struct{}

func (UnimplementedBettingServiceServer) CreateBetting(context.Context, *BettingMessageRequest) (*BettingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBetting not implemented")
}
func (UnimplementedBettingServiceServer) mustEmbedUnimplementedBettingServiceServer() {}
func (UnimplementedBettingServiceServer) testEmbeddedByValue()                        {}

// UnsafeBettingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BettingServiceServer will
// result in compilation errors.
type UnsafeBettingServiceServer interface {
	mustEmbedUnimplementedBettingServiceServer()
}

func RegisterBettingServiceServer(s grpc.ServiceRegistrar, srv BettingServiceServer) {
	// If the following call pancis, it indicates UnimplementedBettingServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BettingService_ServiceDesc, srv)
}

func _BettingService_CreateBetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BettingMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BettingServiceServer).CreateBetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BettingService_CreateBetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BettingServiceServer).CreateBetting(ctx, req.(*BettingMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BettingService_ServiceDesc is the grpc.ServiceDesc for BettingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BettingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "betting.BettingService",
	HandlerType: (*BettingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBetting",
			Handler:    _BettingService_CreateBetting_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/betting/betting.proto",
}
