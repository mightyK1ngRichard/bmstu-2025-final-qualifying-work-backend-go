// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: profile.proto

package generated

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ProfileService_GetUserInfo_FullMethodName         = "/profile.ProfileService/GetUserInfo"
	ProfileService_GetUserInfoByID_FullMethodName     = "/profile.ProfileService/GetUserInfoByID"
	ProfileService_GetUserAddresses_FullMethodName    = "/profile.ProfileService/GetUserAddresses"
	ProfileService_UpdateUserAddresses_FullMethodName = "/profile.ProfileService/UpdateUserAddresses"
	ProfileService_CreateAddress_FullMethodName       = "/profile.ProfileService/CreateAddress"
)

// ProfileServiceClient is the client API for ProfileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// ############### ProfileService ###############
type ProfileServiceClient interface {
	GetUserInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUserInfoRes, error)
	GetUserInfoByID(ctx context.Context, in *GetUserInfoByIDReq, opts ...grpc.CallOption) (*GetUserInfoByIDRes, error)
	GetUserAddresses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUserAddressesRes, error)
	UpdateUserAddresses(ctx context.Context, in *UpdateUserAddressesReq, opts ...grpc.CallOption) (*UpdateUserAddressesRes, error)
	CreateAddress(ctx context.Context, in *CreateAddressReq, opts ...grpc.CallOption) (*CreateAddressRes, error)
}

type profileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileServiceClient(cc grpc.ClientConnInterface) ProfileServiceClient {
	return &profileServiceClient{cc}
}

func (c *profileServiceClient) GetUserInfo(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUserInfoRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserInfoRes)
	err := c.cc.Invoke(ctx, ProfileService_GetUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetUserInfoByID(ctx context.Context, in *GetUserInfoByIDReq, opts ...grpc.CallOption) (*GetUserInfoByIDRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserInfoByIDRes)
	err := c.cc.Invoke(ctx, ProfileService_GetUserInfoByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetUserAddresses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUserAddressesRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserAddressesRes)
	err := c.cc.Invoke(ctx, ProfileService_GetUserAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) UpdateUserAddresses(ctx context.Context, in *UpdateUserAddressesReq, opts ...grpc.CallOption) (*UpdateUserAddressesRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserAddressesRes)
	err := c.cc.Invoke(ctx, ProfileService_UpdateUserAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) CreateAddress(ctx context.Context, in *CreateAddressReq, opts ...grpc.CallOption) (*CreateAddressRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAddressRes)
	err := c.cc.Invoke(ctx, ProfileService_CreateAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServiceServer is the server API for ProfileService service.
// All implementations must embed UnimplementedProfileServiceServer
// for forward compatibility.
//
// ############### ProfileService ###############
type ProfileServiceServer interface {
	GetUserInfo(context.Context, *emptypb.Empty) (*GetUserInfoRes, error)
	GetUserInfoByID(context.Context, *GetUserInfoByIDReq) (*GetUserInfoByIDRes, error)
	GetUserAddresses(context.Context, *emptypb.Empty) (*GetUserAddressesRes, error)
	UpdateUserAddresses(context.Context, *UpdateUserAddressesReq) (*UpdateUserAddressesRes, error)
	CreateAddress(context.Context, *CreateAddressReq) (*CreateAddressRes, error)
	mustEmbedUnimplementedProfileServiceServer()
}

// UnimplementedProfileServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedProfileServiceServer struct{}

func (UnimplementedProfileServiceServer) GetUserInfo(context.Context, *emptypb.Empty) (*GetUserInfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedProfileServiceServer) GetUserInfoByID(context.Context, *GetUserInfoByIDReq) (*GetUserInfoByIDRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoByID not implemented")
}
func (UnimplementedProfileServiceServer) GetUserAddresses(context.Context, *emptypb.Empty) (*GetUserAddressesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserAddresses not implemented")
}
func (UnimplementedProfileServiceServer) UpdateUserAddresses(context.Context, *UpdateUserAddressesReq) (*UpdateUserAddressesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserAddresses not implemented")
}
func (UnimplementedProfileServiceServer) CreateAddress(context.Context, *CreateAddressReq) (*CreateAddressRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAddress not implemented")
}
func (UnimplementedProfileServiceServer) mustEmbedUnimplementedProfileServiceServer() {}
func (UnimplementedProfileServiceServer) testEmbeddedByValue()                        {}

// UnsafeProfileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileServiceServer will
// result in compilation errors.
type UnsafeProfileServiceServer interface {
	mustEmbedUnimplementedProfileServiceServer()
}

func RegisterProfileServiceServer(s grpc.ServiceRegistrar, srv ProfileServiceServer) {
	// If the following call pancis, it indicates UnimplementedProfileServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProfileService_ServiceDesc, srv)
}

func _ProfileService_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetUserInfo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetUserInfoByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoByIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetUserInfoByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetUserInfoByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetUserInfoByID(ctx, req.(*GetUserInfoByIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetUserAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetUserAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetUserAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetUserAddresses(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_UpdateUserAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserAddressesReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).UpdateUserAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_UpdateUserAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).UpdateUserAddresses(ctx, req.(*UpdateUserAddressesReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_CreateAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAddressReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).CreateAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_CreateAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).CreateAddress(ctx, req.(*CreateAddressReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ProfileService_ServiceDesc is the grpc.ServiceDesc for ProfileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profile.ProfileService",
	HandlerType: (*ProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserInfo",
			Handler:    _ProfileService_GetUserInfo_Handler,
		},
		{
			MethodName: "GetUserInfoByID",
			Handler:    _ProfileService_GetUserInfoByID_Handler,
		},
		{
			MethodName: "GetUserAddresses",
			Handler:    _ProfileService_GetUserAddresses_Handler,
		},
		{
			MethodName: "UpdateUserAddresses",
			Handler:    _ProfileService_UpdateUserAddresses_Handler,
		},
		{
			MethodName: "CreateAddress",
			Handler:    _ProfileService_CreateAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}
