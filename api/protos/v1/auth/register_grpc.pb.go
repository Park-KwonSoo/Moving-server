// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: api/protos/v1/auth/register.proto

package auth

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

// RegisterServiceClient is the client API for RegisterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegisterServiceClient interface {
	//회원가입 api
	Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error)
}

type registerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRegisterServiceClient(cc grpc.ClientConnInterface) RegisterServiceClient {
	return &registerServiceClient{cc}
}

func (c *registerServiceClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterRes, error) {
	out := new(RegisterRes)
	err := c.cc.Invoke(ctx, "/v1.register_proto.RegisterService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegisterServiceServer is the server API for RegisterService service.
// All implementations must embed UnimplementedRegisterServiceServer
// for forward compatibility
type RegisterServiceServer interface {
	//회원가입 api
	Register(context.Context, *RegisterReq) (*RegisterRes, error)
	mustEmbedUnimplementedRegisterServiceServer()
}

// UnimplementedRegisterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRegisterServiceServer struct {
}

func (UnimplementedRegisterServiceServer) Register(context.Context, *RegisterReq) (*RegisterRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRegisterServiceServer) mustEmbedUnimplementedRegisterServiceServer() {}

// UnsafeRegisterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegisterServiceServer will
// result in compilation errors.
type UnsafeRegisterServiceServer interface {
	mustEmbedUnimplementedRegisterServiceServer()
}

func RegisterRegisterServiceServer(s grpc.ServiceRegistrar, srv RegisterServiceServer) {
	s.RegisterService(&RegisterService_ServiceDesc, srv)
}

func _RegisterService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegisterServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.register_proto.RegisterService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegisterServiceServer).Register(ctx, req.(*RegisterReq))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterService_ServiceDesc is the grpc.ServiceDesc for RegisterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegisterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.register_proto.RegisterService",
	HandlerType: (*RegisterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _RegisterService_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/protos/v1/auth/register.proto",
}
