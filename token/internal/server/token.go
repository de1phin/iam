package server

import (
	"context"
	"net"
	"sync"

	token "github.com/de1phin/iam/genproto/services/token/api"
	"github.com/de1phin/iam/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StartTokenService(ctx context.Context, serv token.TokenServiceServer, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", ":8080") // TODO отбибабобьте меня и поправьте адрес
	if err != nil {
		logger.Error("failed to listen in sender server", zap.Error(err))
	}
	server := grpc.NewServer()
	token.RegisterTokenServiceServer(server, serv)

	go func() {
		logger.Info("sender server is begin")
		defer wg.Done()

		go func() {
			if err := server.Serve(listener); err != nil {
				logger.Error("failed to serve in sender server", zap.Error(err))
			}
		}()

		<-ctx.Done()
		server.GracefulStop()

		logger.Info("sender server is end")
	}()
	return
}
