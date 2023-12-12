package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	yc "github.com/de1phin/iam/iamssh/pkg/yccli"
	pkgdatabase "github.com/de1phin/iam/pkg/database"
	"github.com/de1phin/iam/pkg/logger"
	"github.com/de1phin/iam/services/account/internal/database"
	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/database/sql"
	"github.com/de1phin/iam/services/account/internal/server"
	"github.com/de1phin/iam/services/account/internal/service"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server                          *server.AccountServiceServerConfig `yaml:"server"`
	Postgres                        *Postgres                          `yaml:"postgres"`
	PostgresPasswordLockboxSecretID string                             `yaml:"postgres_password_lockbox_secret_id"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbname"`
	Password string
	Dsn      string
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

func getPostgresDsn(cfg *Config) string {
	passwordSecret, err := yc.NewLockboxClient().WithTokenGetter(yc.ComputeMetadata()).GetSecret(cfg.PostgresPasswordLockboxSecretID)
	if err != nil {
		logger.Fatal("Failed to get secret from lockbox", zap.Error(err))
	}

	password := passwordSecret.Entries[0].TextValue
	cfg.Postgres.Password = password

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
}

func main() {
	configPath := flag.String("config", "", "config path")
	flag.Parse()

	config := readConfig(*configPath)

	var accountDb database.AccountDatabase
	var sshKeyDb database.SshKeyDatabase

	if config.Postgres == nil {
		accountDb = cache.NewAccountCache()
		sshKeyDb = cache.NewSshKeyCache()
	} else {
		dsn := getPostgresDsn(config)
		database, err := pkgdatabase.NewDatabase(context.Background(), dsn)
		if err != nil {
			log.Fatal(err)
		}
		sqlDatabase, err := sql.New(context.Background(), database)
		if err != nil {
			log.Fatal(err)
		}
		accountDb = sql.NewSqlAccountDatabase(sqlDatabase)
		sshKeyDb = sql.NewSqlSshKeyDatabase(sqlDatabase)
	}

	service := service.NewAccountService(
		service.AccountDatabase(accountDb),
		service.SshKeyDatabase(sshKeyDb),
	)

	zapCfg := zap.NewProductionConfig()
	zapCfg.OutputPaths = []string{"stdout"}
	logger := zap.Must(zapCfg.Build()).Named("AccountService")

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
