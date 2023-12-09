package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/server"
	"github.com/de1phin/iam/services/account/internal/service"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server *server.AccountServiceServerConfig `yaml:"server"`
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

	serverConfig := config.Server
	serverConfig.AccountService = service
	serverConfig.Logger = logger

	srv, err := serverConfig.RunServer()
	if err != nil {
		logger.Fatal("Failed to run server", zap.Error(err))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGTERM)

	signal := <-ch
	logger.Info("Exiting on signal", zap.Any("signal", signal))
	srv.Shutdown()
}
