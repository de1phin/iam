// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: access/api/access-service.proto

package token

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

const (
	AccessService_AddRole_FullMethodName             = "/iam.services.auth.AccessService/AddRole"
	AccessService_GetRole_FullMethodName             = "/iam.services.auth.AccessService/GetRole"
	AccessService_DeleteRole_FullMethodName          = "/iam.services.auth.AccessService/DeleteRole"
	AccessService_AddAccessBinding_FullMethodName    = "/iam.services.auth.AccessService/AddAccessBinding"
	AccessService_CheckPermission_FullMethodName     = "/iam.services.auth.AccessService/CheckPermission"
	AccessService_DeleteAccessBinding_FullMethodName = "/iam.services.auth.AccessService/DeleteAccessBinding"
)

// AccessServiceClient is the client API for AccessService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessServiceClient interface {
	AddRole(ctx context.Context, in *AddRoleRequest, opts ...grpc.CallOption) (*AddRoleResponse, error)
	GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*GetRoleResponse, error)
	DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error)
	AddAccessBinding(ctx context.Context, in *AddAccessBindingRequest, opts ...grpc.CallOption) (*AddAccessBindingResponse, error)
	CheckPermission(ctx context.Context, in *CheckPermissionRequest, opts ...grpc.CallOption) (*CheckPermissionResponse, error)
	DeleteAccessBinding(ctx context.Context, in *DeleteAccessBindingRequest, opts ...grpc.CallOption) (*DeleteAccessBindingResponse, error)
}

type accessServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccessServiceClient(cc grpc.ClientConnInterface) AccessServiceClient {
	return &accessServiceClient{cc}
}

func (c *accessServiceClient) AddRole(ctx context.Context, in *AddRoleRequest, opts ...grpc.CallOption) (*AddRoleResponse, error) {
	out := new(AddRoleResponse)
	err := c.cc.Invoke(ctx, AccessService_AddRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) GetRole(ctx context.Context, in *GetRoleRequest, opts ...grpc.CallOption) (*GetRoleResponse, error) {
	out := new(GetRoleResponse)
	err := c.cc.Invoke(ctx, AccessService_GetRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...grpc.CallOption) (*DeleteRoleResponse, error) {
	out := new(DeleteRoleResponse)
	err := c.cc.Invoke(ctx, AccessService_DeleteRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) AddAccessBinding(ctx context.Context, in *AddAccessBindingRequest, opts ...grpc.CallOption) (*AddAccessBindingResponse, error) {
	out := new(AddAccessBindingResponse)
	err := c.cc.Invoke(ctx, AccessService_AddAccessBinding_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) CheckPermission(ctx context.Context, in *CheckPermissionRequest, opts ...grpc.CallOption) (*CheckPermissionResponse, error) {
	out := new(CheckPermissionResponse)
	err := c.cc.Invoke(ctx, AccessService_CheckPermission_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessServiceClient) DeleteAccessBinding(ctx context.Context, in *DeleteAccessBindingRequest, opts ...grpc.CallOption) (*DeleteAccessBindingResponse, error) {
	out := new(DeleteAccessBindingResponse)
	err := c.cc.Invoke(ctx, AccessService_DeleteAccessBinding_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessServiceServer is the server API for AccessService service.
// All implementations must embed UnimplementedAccessServiceServer
// for forward compatibility
type AccessServiceServer interface {
	AddRole(context.Context, *AddRoleRequest) (*AddRoleResponse, error)
	GetRole(context.Context, *GetRoleRequest) (*GetRoleResponse, error)
	DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error)
	AddAccessBinding(context.Context, *AddAccessBindingRequest) (*AddAccessBindingResponse, error)
	CheckPermission(context.Context, *CheckPermissionRequest) (*CheckPermissionResponse, error)
	DeleteAccessBinding(context.Context, *DeleteAccessBindingRequest) (*DeleteAccessBindingResponse, error)
	mustEmbedUnimplementedAccessServiceServer()
}

// UnimplementedAccessServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccessServiceServer struct {
}

func (UnimplementedAccessServiceServer) AddRole(context.Context, *AddRoleRequest) (*AddRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRole not implemented")
}
func (UnimplementedAccessServiceServer) GetRole(context.Context, *GetRoleRequest) (*GetRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRole not implemented")
}
func (UnimplementedAccessServiceServer) DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRole not implemented")
}
func (UnimplementedAccessServiceServer) AddAccessBinding(context.Context, *AddAccessBindingRequest) (*AddAccessBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAccessBinding not implemented")
}
func (UnimplementedAccessServiceServer) CheckPermission(context.Context, *CheckPermissionRequest) (*CheckPermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPermission not implemented")
}
func (UnimplementedAccessServiceServer) DeleteAccessBinding(context.Context, *DeleteAccessBindingRequest) (*DeleteAccessBindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccessBinding not implemented")
}
func (UnimplementedAccessServiceServer) mustEmbedUnimplementedAccessServiceServer() {}

// UnsafeAccessServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessServiceServer will
// result in compilation errors.
type UnsafeAccessServiceServer interface {
	mustEmbedUnimplementedAccessServiceServer()
}

func RegisterAccessServiceServer(s grpc.ServiceRegistrar, srv AccessServiceServer) {
	s.RegisterService(&AccessService_ServiceDesc, srv)
}

func _AccessService_AddRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).AddRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessService_AddRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).AddRole(ctx, req.(*AddRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_GetRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).GetRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessService_GetRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).GetRole(ctx, req.(*GetRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_DeleteRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).DeleteRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessService_DeleteRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).DeleteRole(ctx, req.(*DeleteRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_AddAccessBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAccessBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).AddAccessBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessService_AddAccessBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).AddAccessBinding(ctx, req.(*AddAccessBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_CheckPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).CheckPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessService_CheckPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).CheckPermission(ctx, req.(*CheckPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessService_DeleteAccessBinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccessBindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessServiceServer).DeleteAccessBinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccessService_DeleteAccessBinding_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessServiceServer).DeleteAccessBinding(ctx, req.(*DeleteAccessBindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccessService_ServiceDesc is the grpc.ServiceDesc for AccessService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccessService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "iam.services.auth.AccessService",
	HandlerType: (*AccessServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddRole",
			Handler:    _AccessService_AddRole_Handler,
		},
		{
			MethodName: "GetRole",
			Handler:    _AccessService_GetRole_Handler,
		},
		{
			MethodName: "DeleteRole",
			Handler:    _AccessService_DeleteRole_Handler,
		},
		{
			MethodName: "AddAccessBinding",
			Handler:    _AccessService_AddAccessBinding_Handler,
		},
		{
			MethodName: "CheckPermission",
			Handler:    _AccessService_CheckPermission_Handler,
		},
		{
			MethodName: "DeleteAccessBinding",
			Handler:    _AccessService_DeleteAccessBinding_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "access/api/access-service.proto",
}