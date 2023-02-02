// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: proto.proto

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RegClient is the client API for Reg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegClient interface {
	// for register service implimentation
	SignUp(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type regClient struct {
	cc grpc.ClientConnInterface
}

func NewRegClient(cc grpc.ClientConnInterface) RegClient {
	return &regClient{cc}
}

func (c *regClient) SignUp(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/protobuf.Reg/SignUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegServer is the server API for Reg service.
// All implementations must embed UnimplementedRegServer
// for forward compatibility
type RegServer interface {
	// for register service implimentation
	SignUp(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedRegServer()
}

// UnimplementedRegServer must be embedded to have forward compatible implementations.
type UnimplementedRegServer struct {
}

func (UnimplementedRegServer) SignUp(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUp not implemented")
}
func (UnimplementedRegServer) mustEmbedUnimplementedRegServer() {}

// UnsafeRegServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegServer will
// result in compilation errors.
type UnsafeRegServer interface {
	mustEmbedUnimplementedRegServer()
}

func RegisterRegServer(s grpc.ServiceRegistrar, srv RegServer) {
	s.RegisterService(&Reg_ServiceDesc, srv)
}

func _Reg_SignUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegServer).SignUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.Reg/SignUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegServer).SignUp(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Reg_ServiceDesc is the grpc.ServiceDesc for Reg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Reg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.Reg",
	HandlerType: (*RegServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignUp",
			Handler:    _Reg_SignUp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto.proto",
}
