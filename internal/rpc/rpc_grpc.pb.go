// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: internal/rpc/rpc.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Gokeeper_Register_FullMethodName       = "/rpc.Gokeeper/Register"
	Gokeeper_Login_FullMethodName          = "/rpc.Gokeeper/Login"
	Gokeeper_GetRecordsInfo_FullMethodName = "/rpc.Gokeeper/GetRecordsInfo"
	Gokeeper_GetRecord_FullMethodName      = "/rpc.Gokeeper/GetRecord"
	Gokeeper_CreateRecord_FullMethodName   = "/rpc.Gokeeper/CreateRecord"
	Gokeeper_DeleteRecord_FullMethodName   = "/rpc.Gokeeper/DeleteRecord"
)

// GokeeperClient is the client API for Gokeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GokeeperClient interface {
	Register(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Session, error)
	Login(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Session, error)
	GetRecordsInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RecordsList, error)
	GetRecord(ctx context.Context, in *RecordID, opts ...grpc.CallOption) (*Record, error)
	CreateRecord(ctx context.Context, in *Record, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteRecord(ctx context.Context, in *RecordID, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type gokeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGokeeperClient(cc grpc.ClientConnInterface) GokeeperClient {
	return &gokeeperClient{cc}
}

func (c *gokeeperClient) Register(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, Gokeeper_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gokeeperClient) Login(ctx context.Context, in *UserCredentials, opts ...grpc.CallOption) (*Session, error) {
	out := new(Session)
	err := c.cc.Invoke(ctx, Gokeeper_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gokeeperClient) GetRecordsInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*RecordsList, error) {
	out := new(RecordsList)
	err := c.cc.Invoke(ctx, Gokeeper_GetRecordsInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gokeeperClient) GetRecord(ctx context.Context, in *RecordID, opts ...grpc.CallOption) (*Record, error) {
	out := new(Record)
	err := c.cc.Invoke(ctx, Gokeeper_GetRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gokeeperClient) CreateRecord(ctx context.Context, in *Record, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gokeeper_CreateRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gokeeperClient) DeleteRecord(ctx context.Context, in *RecordID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Gokeeper_DeleteRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GokeeperServer is the server API for Gokeeper service.
// All implementations must embed UnimplementedGokeeperServer
// for forward compatibility
type GokeeperServer interface {
	Register(context.Context, *UserCredentials) (*Session, error)
	Login(context.Context, *UserCredentials) (*Session, error)
	GetRecordsInfo(context.Context, *emptypb.Empty) (*RecordsList, error)
	GetRecord(context.Context, *RecordID) (*Record, error)
	CreateRecord(context.Context, *Record) (*emptypb.Empty, error)
	DeleteRecord(context.Context, *RecordID) (*emptypb.Empty, error)
	mustEmbedUnimplementedGokeeperServer()
}

// UnimplementedGokeeperServer must be embedded to have forward compatible implementations.
type UnimplementedGokeeperServer struct {
}

func (UnimplementedGokeeperServer) Register(context.Context, *UserCredentials) (*Session, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedGokeeperServer) Login(context.Context, *UserCredentials) (*Session, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedGokeeperServer) GetRecordsInfo(context.Context, *emptypb.Empty) (*RecordsList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecordsInfo not implemented")
}
func (UnimplementedGokeeperServer) GetRecord(context.Context, *RecordID) (*Record, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecord not implemented")
}
func (UnimplementedGokeeperServer) CreateRecord(context.Context, *Record) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRecord not implemented")
}
func (UnimplementedGokeeperServer) DeleteRecord(context.Context, *RecordID) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRecord not implemented")
}
func (UnimplementedGokeeperServer) mustEmbedUnimplementedGokeeperServer() {}

// UnsafeGokeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GokeeperServer will
// result in compilation errors.
type UnsafeGokeeperServer interface {
	mustEmbedUnimplementedGokeeperServer()
}

func RegisterGokeeperServer(s grpc.ServiceRegistrar, srv GokeeperServer) {
	s.RegisterService(&Gokeeper_ServiceDesc, srv)
}

func _Gokeeper_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCredentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GokeeperServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gokeeper_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GokeeperServer).Register(ctx, req.(*UserCredentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gokeeper_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCredentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GokeeperServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gokeeper_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GokeeperServer).Login(ctx, req.(*UserCredentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gokeeper_GetRecordsInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GokeeperServer).GetRecordsInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gokeeper_GetRecordsInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GokeeperServer).GetRecordsInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gokeeper_GetRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GokeeperServer).GetRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gokeeper_GetRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GokeeperServer).GetRecord(ctx, req.(*RecordID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gokeeper_CreateRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Record)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GokeeperServer).CreateRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gokeeper_CreateRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GokeeperServer).CreateRecord(ctx, req.(*Record))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gokeeper_DeleteRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GokeeperServer).DeleteRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Gokeeper_DeleteRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GokeeperServer).DeleteRecord(ctx, req.(*RecordID))
	}
	return interceptor(ctx, in, info, handler)
}

// Gokeeper_ServiceDesc is the grpc.ServiceDesc for Gokeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Gokeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Gokeeper",
	HandlerType: (*GokeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Gokeeper_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Gokeeper_Login_Handler,
		},
		{
			MethodName: "GetRecordsInfo",
			Handler:    _Gokeeper_GetRecordsInfo_Handler,
		},
		{
			MethodName: "GetRecord",
			Handler:    _Gokeeper_GetRecord_Handler,
		},
		{
			MethodName: "CreateRecord",
			Handler:    _Gokeeper_CreateRecord_Handler,
		},
		{
			MethodName: "DeleteRecord",
			Handler:    _Gokeeper_DeleteRecord_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/rpc/rpc.proto",
}
