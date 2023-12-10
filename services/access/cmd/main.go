package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	yc "github.com/de1phin/iam/iamssh/pkg/yccli"
	"github.com/de1phin/iam/services/access/internal/grpcclient"
	"github.com/de1phin/iam/services/access/internal/postgres"
	"github.com/de1phin/iam/services/access/internal/server"
	"github.com/de1phin/iam/services/access/internal/service"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server                          *server.AccessServiceServerConfig `yaml:"server"`
	TokenServiceAddress             string                            `yaml:"token_service_addr"`
	Postgres                        *postgres.Options                 `yaml:"postgres"`
	PostgresPasswordLockboxSecretID string                            `yaml:"postgres_password_lockbox_secret_id"`
}

func readConfig(configPath string) *Config {
	var cfg Config
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("failed to read config from ", configPath, ":", err)
	}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		log.Fatal("failed to unmarshal config - ", err)
	}
	return &cfg
}

func main() {
	configPath := flag.String("config", "", "config path")
	flag.Parse()

	config := readConfig(*configPath)

	ctx := context.Background()

	zapCfg := zap.NewProductionConfig()
	zapCfg.OutputPaths = []string{"stdout"}
	logger := zap.Must(zapCfg.Build()).Named("AccessService")

	tokenValidator, err := grpcclient.NewTokenServiceClient(ctx, config.TokenServiceAddress)
	if err != nil {
		logger.Fatal("Failed to establish connection to token service", zap.Error(err))
	}

	passwordSecret, err := yc.NewLockboxClient().WithTokenGetter(yc.ComputeMetadata()).GetSecret(config.PostgresPasswordLockboxSecretID)
	if err != nil {
		logger.Fatal("Failed to get secret from lockbox", zap.Error(err))
	}

	password := passwordSecret.Entries[0].TextValue
	config.Postgres.Password = password

	storage, err := postgres.New(ctx, *config.Postgres)
	if err != nil {
		logger.Fatal("Failed to establish connection to postgresql", zap.Error(err))
	}

	service := service.New(storage, tokenValidator)

	serverConfig := config.Server

	srv, err := serverConfig.RunServer(service, logger)
	if err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTERM)

	signal := <-ch
	logger.Info("Exiting on signal", zap.Any("signal", signal))
	srv.Shutdown()
}
