package server

import (
	"context"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	token "github.com/de1phin/iam/genproto/services/token/api"
	"github.com/de1phin/iam/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func InitTokenSwagger(ctx context.Context, wg *sync.WaitGroup, swaggerHost, grpcHost string) {
	httpMux := http.NewServeMux()

	relativePath := "./genproto/services/token/api/token-service.swagger.json"
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		logger.Error("build absolutePath", zap.Error(err))
	}

	httpMux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, absolutePath)
	})

	httpMux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://"+swaggerHost+"/swagger.json"),
	))

	grpcMux := runtime.NewServeMux()
	if err := token.RegisterTokenServiceHandlerFromEndpoint(ctx, grpcMux, grpcHost, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		logger.Error("failed to register gateway handler", zap.Error(err))
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		srv := &http.Server{
			Addr: swaggerHost,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/swagger") {
					httpMux.ServeHTTP(w, r)
					return
				}
				grpcMux.ServeHTTP(w, r)
			}),
		}

		logger.Info("swagger for token-service start")
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				logger.Error("swagger for token-service", zap.Error(err))
			}
		}()

		<-ctx.Done()
		_ = srv.Shutdown(ctx)

		logger.Info("swagger for token-service stop")
	}()
}

func StartTokenService(ctx context.Context, serv token.TokenServiceServer, wg *sync.WaitGroup, host string) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		logger.Error("failed to listen in sender server", zap.Error(err))
	}
	server := grpc.NewServer()
	token.RegisterTokenServiceServer(server, serv)
	reflection.Register(server)

	wg.Add(1)
	go func() {
		logger.Info("token-service start")
		defer wg.Done()

		go func() {
			if err := server.Serve(listener); err != nil {
				logger.Error("failed to serve in sender server", zap.Error(err))
			}
		}()

		<-ctx.Done()
		server.GracefulStop()

		logger.Info("sender server stop")
	}()

}
