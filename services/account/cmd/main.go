package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/server"
	"github.com/de1phin/iam/services/account/internal/service"
	"go.uber.org/zap"
)

func main() {

	accountCache := cache.NewAccountCache()
	sshKeysCache := cache.NewSshKeyCache()

	service := service.NewAccountService(
		service.AccountDatabase(accountCache),
		service.SshKeyDatabase(sshKeysCache),
	)

	logger := zap.Must(zap.NewProduction()).Named("AccountService")
	instance := os.Getenv("INSTANCE")
	if instance != "" {
		logger = logger.With(zap.String("instance", instance))
	}

	cfg := &server.AccountServiceServerConfig{
		AccountService:    service,
		Logger:            logger,
		Address:           ":8443",
		ConnectionTimeout: time.Second * 30,
		// *tls.Certificate
	}

	srv, err := cfg.RunServer()
	if err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTERM)

	signal := <-ch
	logger.Info("Exiting on signal", zap.Any("signal", signal))
	srv.Shutdown()
}
