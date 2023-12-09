package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type AccountServiceServerConfig struct {
	AccountService    *service.AccountService
	Logger            *zap.Logger
	Address           string        `yaml:"address"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
	TlsCertificate    *tls.Certificate
}

func (cfg *AccountServiceServerConfig) getGrpcServerOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}

	if cfg.ConnectionTimeout != 0 {
		opts = append(opts, grpc.ConnectionTimeout(cfg.ConnectionTimeout))
	}

	if cfg.TlsCertificate != nil {
		c := credentials.NewServerTLSFromCert(cfg.TlsCertificate)
		opts = append(opts, grpc.Creds(c))
	}

	return opts
}

type Server struct {
	gprcServer *grpc.Server
	logger     *zap.Logger
}

func (cfg *AccountServiceServerConfig) RunServer() (*Server, error) {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen tcp: %w", err)
	}

	grpcServer := grpc.NewServer(cfg.getGrpcServerOptions()...)
	account.RegisterAccountServiceServer(grpcServer, cfg.AccountService)
	reflection.Register(grpcServer)

	srv := &Server{
		gprcServer: grpcServer,
		logger:     cfg.Logger,
	}

	go func() {
		srv.logger.Info("Run account service server")
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return srv, nil
}

func (s *Server) Shutdown() {
	s.logger.Info("Gracefully shutdown account service server")
	s.gprcServer.GracefulStop()
}
