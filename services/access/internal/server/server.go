package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	access "github.com/de1phin/iam/genproto/services/access/api"
	"github.com/de1phin/iam/services/access/internal/core"
	"github.com/de1phin/iam/services/access/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	AccessService *service.AccessService

	*access.UnimplementedAccessServiceServer

	gprcServer *grpc.Server
	logger     *zap.Logger
}

type AccessServiceServerConfig struct {
	Address           string        `yaml:"address"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
}

func (cfg *AccessServiceServerConfig) getGrpcServerOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}

	if cfg.ConnectionTimeout != 0 {
		opts = append(opts, grpc.ConnectionTimeout(cfg.ConnectionTimeout))
	}

	return opts
}

func (cfg *AccessServiceServerConfig) RunServer(srv *service.AccessService, logger *zap.Logger) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen tcp: %w", err)
	}

	server := &Server{
		logger:        logger,
		AccessService: srv,
	}

	grpcServer := grpc.NewServer(cfg.getGrpcServerOptions()...)
	access.RegisterAccessServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	server.gprcServer = grpcServer

	go func() {
		server.logger.Info("Run access service server")
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return server, nil
}

func (s *Server) Shutdown() {
	s.logger.Info("Gracefully shutdown access service server")
	s.gprcServer.GracefulStop()
}

func (s *Server) AddRole(ctx context.Context, req *access.AddRoleRequest) (*access.AddRoleResponse, error) {
	err := s.AccessService.AddRole(ctx, core.Role{
		Name:        req.GetRole().GetName(),
		Permissions: req.GetRole().GetPermissions(),
	})

	if err != nil {
		return nil, InternalError(err)
	}

	return &access.AddRoleResponse{}, nil
}

func (s *Server) GetRole(ctx context.Context, req *access.GetRoleRequest) (*access.GetRoleResponse, error) {
	role, err := s.AccessService.GetRole(ctx, req.GetName())
	if err != nil {
		return nil, InternalError(err)
	}

	return &access.GetRoleResponse{
		Role: &access.Role{
			Name:        role.Name,
			Permissions: role.Permissions,
		},
	}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *access.DeleteRoleRequest) (*access.DeleteRoleResponse, error) {
	err := s.AccessService.DeleteRole(ctx, req.GetName())
	if err != nil {
		return nil, InternalError(err)
	}

	return &access.DeleteRoleResponse{}, nil
}

func (s *Server) AddAccessBinding(ctx context.Context, req *access.AddAccessBindingRequest) (*access.AddAccessBindingResponse, error) {
	err := s.AccessService.AddAccessBinding(ctx, core.AccessBinding{
		AccountID: req.GetAccessBinding().GetAccountId(),
		RoleName:  req.GetAccessBinding().GetRoleName(),
		Resource:  req.GetAccessBinding().GetResource(),
	})

	if err != nil {
		return nil, InternalError(err)
	}

	return &access.AddAccessBindingResponse{}, nil
}

func (s *Server) CheckPermission(ctx context.Context, req *access.CheckPermissionRequest) (*access.CheckPermissionResponse, error) {
	ok, err := s.AccessService.CheckPermission(ctx,
		req.GetToken(), req.GetResource(), req.GetPermission())
	if errors.Is(err, core.ErrInvalidToken) {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated (Invalid token)")
	}
	if err != nil {
		return nil, InternalError(err)
	}

	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	return &access.CheckPermissionResponse{}, nil
}

func (s *Server) DeleteAccessBinding(ctx context.Context, req *access.DeleteAccessBindingRequest) (*access.DeleteAccessBindingResponse, error) {
	err := s.AccessService.DeleteAccessBinding(ctx, core.AccessBinding{
		AccountID: req.GetAccessBinding().GetAccountId(),
		RoleName:  req.GetAccessBinding().GetRoleName(),
		Resource:  req.GetAccessBinding().GetResource(),
	})

	if err != nil {
		return nil, InternalError(err)
	}

	return &access.DeleteAccessBindingResponse{}, nil
}
