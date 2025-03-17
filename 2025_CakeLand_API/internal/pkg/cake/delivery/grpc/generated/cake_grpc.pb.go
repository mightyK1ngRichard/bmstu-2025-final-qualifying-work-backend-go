// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: cake.proto

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
	CakeService_Cake_FullMethodName           = "/CakeService/Cake"
	CakeService_CreateCake_FullMethodName     = "/CakeService/CreateCake"
	CakeService_CreateFilling_FullMethodName  = "/CakeService/CreateFilling"
	CakeService_CreateCategory_FullMethodName = "/CakeService/CreateCategory"
	CakeService_Categories_FullMethodName     = "/CakeService/Categories"
	CakeService_Fillings_FullMethodName       = "/CakeService/Fillings"
	CakeService_Cakes_FullMethodName          = "/CakeService/Cakes"
)

// CakeServiceClient is the client API for CakeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CakeServiceClient interface {
	Cake(ctx context.Context, in *CakeRequest, opts ...grpc.CallOption) (*CakeResponse, error)
	CreateCake(ctx context.Context, in *CreateCakeRequest, opts ...grpc.CallOption) (*CreateCakeResponse, error)
	CreateFilling(ctx context.Context, in *CreateFillingRequest, opts ...grpc.CallOption) (*CreateFillingResponse, error)
	CreateCategory(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*CreateCategoryResponse, error)
	Categories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CategoriesResponse, error)
	Fillings(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FillingsResponse, error)
	Cakes(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CakesResponse, error)
}

type cakeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCakeServiceClient(cc grpc.ClientConnInterface) CakeServiceClient {
	return &cakeServiceClient{cc}
}

func (c *cakeServiceClient) Cake(ctx context.Context, in *CakeRequest, opts ...grpc.CallOption) (*CakeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CakeResponse)
	err := c.cc.Invoke(ctx, CakeService_Cake_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cakeServiceClient) CreateCake(ctx context.Context, in *CreateCakeRequest, opts ...grpc.CallOption) (*CreateCakeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCakeResponse)
	err := c.cc.Invoke(ctx, CakeService_CreateCake_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cakeServiceClient) CreateFilling(ctx context.Context, in *CreateFillingRequest, opts ...grpc.CallOption) (*CreateFillingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFillingResponse)
	err := c.cc.Invoke(ctx, CakeService_CreateFilling_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cakeServiceClient) CreateCategory(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*CreateCategoryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCategoryResponse)
	err := c.cc.Invoke(ctx, CakeService_CreateCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cakeServiceClient) Categories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CategoriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CategoriesResponse)
	err := c.cc.Invoke(ctx, CakeService_Categories_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cakeServiceClient) Fillings(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*FillingsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FillingsResponse)
	err := c.cc.Invoke(ctx, CakeService_Fillings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cakeServiceClient) Cakes(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CakesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CakesResponse)
	err := c.cc.Invoke(ctx, CakeService_Cakes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CakeServiceServer is the server API for CakeService service.
// All implementations must embed UnimplementedCakeServiceServer
// for forward compatibility.
type CakeServiceServer interface {
	Cake(context.Context, *CakeRequest) (*CakeResponse, error)
	CreateCake(context.Context, *CreateCakeRequest) (*CreateCakeResponse, error)
	CreateFilling(context.Context, *CreateFillingRequest) (*CreateFillingResponse, error)
	CreateCategory(context.Context, *CreateCategoryRequest) (*CreateCategoryResponse, error)
	Categories(context.Context, *emptypb.Empty) (*CategoriesResponse, error)
	Fillings(context.Context, *emptypb.Empty) (*FillingsResponse, error)
	Cakes(context.Context, *emptypb.Empty) (*CakesResponse, error)
	mustEmbedUnimplementedCakeServiceServer()
}

// UnimplementedCakeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCakeServiceServer struct{}

func (UnimplementedCakeServiceServer) Cake(context.Context, *CakeRequest) (*CakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cake not implemented")
}
func (UnimplementedCakeServiceServer) CreateCake(context.Context, *CreateCakeRequest) (*CreateCakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCake not implemented")
}
func (UnimplementedCakeServiceServer) CreateFilling(context.Context, *CreateFillingRequest) (*CreateFillingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFilling not implemented")
}
func (UnimplementedCakeServiceServer) CreateCategory(context.Context, *CreateCategoryRequest) (*CreateCategoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCategory not implemented")
}
func (UnimplementedCakeServiceServer) Categories(context.Context, *emptypb.Empty) (*CategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Categories not implemented")
}
func (UnimplementedCakeServiceServer) Fillings(context.Context, *emptypb.Empty) (*FillingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fillings not implemented")
}
func (UnimplementedCakeServiceServer) Cakes(context.Context, *emptypb.Empty) (*CakesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Cakes not implemented")
}
func (UnimplementedCakeServiceServer) mustEmbedUnimplementedCakeServiceServer() {}
func (UnimplementedCakeServiceServer) testEmbeddedByValue()                     {}

// UnsafeCakeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CakeServiceServer will
// result in compilation errors.
type UnsafeCakeServiceServer interface {
	mustEmbedUnimplementedCakeServiceServer()
}

func RegisterCakeServiceServer(s grpc.ServiceRegistrar, srv CakeServiceServer) {
	// If the following call pancis, it indicates UnimplementedCakeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CakeService_ServiceDesc, srv)
}

func _CakeService_Cake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakeServiceServer).Cake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CakeService_Cake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakeServiceServer).Cake(ctx, req.(*CakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CakeService_CreateCake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakeServiceServer).CreateCake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CakeService_CreateCake_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakeServiceServer).CreateCake(ctx, req.(*CreateCakeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CakeService_CreateFilling_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFillingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakeServiceServer).CreateFilling(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CakeService_CreateFilling_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakeServiceServer).CreateFilling(ctx, req.(*CreateFillingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CakeService_CreateCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakeServiceServer).CreateCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CakeService_CreateCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakeServiceServer).CreateCategory(ctx, req.(*CreateCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CakeService_Categories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakeServiceServer).Categories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CakeService_Categories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakeServiceServer).Categories(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CakeService_Fillings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakeServiceServer).Fillings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CakeService_Fillings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakeServiceServer).Fillings(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CakeService_Cakes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CakeServiceServer).Cakes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CakeService_Cakes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CakeServiceServer).Cakes(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// CakeService_ServiceDesc is the grpc.ServiceDesc for CakeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CakeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CakeService",
	HandlerType: (*CakeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Cake",
			Handler:    _CakeService_Cake_Handler,
		},
		{
			MethodName: "CreateCake",
			Handler:    _CakeService_CreateCake_Handler,
		},
		{
			MethodName: "CreateFilling",
			Handler:    _CakeService_CreateFilling_Handler,
		},
		{
			MethodName: "CreateCategory",
			Handler:    _CakeService_CreateCategory_Handler,
		},
		{
			MethodName: "Categories",
			Handler:    _CakeService_Categories_Handler,
		},
		{
			MethodName: "Fillings",
			Handler:    _CakeService_Fillings_Handler,
		},
		{
			MethodName: "Cakes",
			Handler:    _CakeService_Cakes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cake.proto",
}
