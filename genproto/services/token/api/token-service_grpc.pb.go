// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: services/token/api/token-service.proto

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
	TokenService_GenerateToken_FullMethodName = "/iam.services.token.TokenService/GenerateToken"
	TokenService_RefreshToken_FullMethodName  = "/iam.services.token.TokenService/RefreshToken"
	TokenService_DeleteToken_FullMethodName   = "/iam.services.token.TokenService/DeleteToken"
	TokenService_CheckToken_FullMethodName    = "/iam.services.token.TokenService/CheckToken"
)

// TokenServiceClient is the client API for TokenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenServiceClient interface {
	GenerateToken(ctx context.Context, in *GenerateTokenRequest, opts ...grpc.CallOption) (*GenerateTokenResponse, error)
	RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error)
	DeleteToken(ctx context.Context, in *RemoveTokenRequest, opts ...grpc.CallOption) (*RemoveTokenResponse, error)
	CheckToken(ctx context.Context, in *CheckTokenRequest, opts ...grpc.CallOption) (*CheckTokenResponse, error)
}

type tokenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenServiceClient(cc grpc.ClientConnInterface) TokenServiceClient {
	return &tokenServiceClient{cc}
}

func (c *tokenServiceClient) GenerateToken(ctx context.Context, in *GenerateTokenRequest, opts ...grpc.CallOption) (*GenerateTokenResponse, error) {
	out := new(GenerateTokenResponse)
	err := c.cc.Invoke(ctx, TokenService_GenerateToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) RefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*RefreshTokenResponse, error) {
	out := new(RefreshTokenResponse)
	err := c.cc.Invoke(ctx, TokenService_RefreshToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) DeleteToken(ctx context.Context, in *RemoveTokenRequest, opts ...grpc.CallOption) (*RemoveTokenResponse, error) {
	out := new(RemoveTokenResponse)
	err := c.cc.Invoke(ctx, TokenService_DeleteToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) CheckToken(ctx context.Context, in *CheckTokenRequest, opts ...grpc.CallOption) (*CheckTokenResponse, error) {
	out := new(CheckTokenResponse)
	err := c.cc.Invoke(ctx, TokenService_CheckToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenServiceServer is the server API for TokenService service.
// All implementations must embed UnimplementedTokenServiceServer
// for forward compatibility
type TokenServiceServer interface {
	GenerateToken(context.Context, *GenerateTokenRequest) (*GenerateTokenResponse, error)
	RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error)
	DeleteToken(context.Context, *RemoveTokenRequest) (*RemoveTokenResponse, error)
	CheckToken(context.Context, *CheckTokenRequest) (*CheckTokenResponse, error)
	mustEmbedUnimplementedTokenServiceServer()
}

// UnimplementedTokenServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTokenServiceServer struct {
}

func (UnimplementedTokenServiceServer) GenerateToken(context.Context, *GenerateTokenRequest) (*GenerateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateToken not implemented")
}
func (UnimplementedTokenServiceServer) RefreshToken(context.Context, *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedTokenServiceServer) DeleteToken(context.Context, *RemoveTokenRequest) (*RemoveTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteToken not implemented")
}
func (UnimplementedTokenServiceServer) CheckToken(context.Context, *CheckTokenRequest) (*CheckTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckToken not implemented")
}
func (UnimplementedTokenServiceServer) mustEmbedUnimplementedTokenServiceServer() {}

// UnsafeTokenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenServiceServer will
// result in compilation errors.
type UnsafeTokenServiceServer interface {
	mustEmbedUnimplementedTokenServiceServer()
}

func RegisterTokenServiceServer(s grpc.ServiceRegistrar, srv TokenServiceServer) {
	s.RegisterService(&TokenService_ServiceDesc, srv)
}

func _TokenService_GenerateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).GenerateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_GenerateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).GenerateToken(ctx, req.(*GenerateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_RefreshToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).RefreshToken(ctx, req.(*RefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_DeleteToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).DeleteToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_DeleteToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).DeleteToken(ctx, req.(*RemoveTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_CheckToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).CheckToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TokenService_CheckToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).CheckToken(ctx, req.(*CheckTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenService_ServiceDesc is the grpc.ServiceDesc for TokenService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "iam.services.token.TokenService",
	HandlerType: (*TokenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateToken",
			Handler:    _TokenService_GenerateToken_Handler,
		},
		{
			MethodName: "RefreshToken",
			Handler:    _TokenService_RefreshToken_Handler,
		},
		{
			MethodName: "DeleteToken",
			Handler:    _TokenService_DeleteToken_Handler,
		},
		{
			MethodName: "CheckToken",
			Handler:    _TokenService_CheckToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/token/api/token-service.proto",
}
