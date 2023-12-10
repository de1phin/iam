package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	account "github.com/de1phin/iam/genproto/services/account/api"
	yc "github.com/de1phin/iam/iamssh/pkg/yccli"
	memcache "github.com/de1phin/iam/pkg/cache"
	"github.com/de1phin/iam/pkg/database"
	"github.com/de1phin/iam/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"

	token_service "github.com/de1phin/iam/services/token/app/token"
	"github.com/de1phin/iam/services/token/internal/cache"
	"github.com/de1phin/iam/services/token/internal/client"
	"github.com/de1phin/iam/services/token/internal/facade"
	"github.com/de1phin/iam/services/token/internal/generator"
	"github.com/de1phin/iam/services/token/internal/repository"
	"github.com/de1phin/iam/services/token/internal/server"
)

type connections struct {
	memcached *memcache.Cache[string, []byte]
	database  database.Database
}

type repositories struct {
	cache *cache.MemCache
	repo  *repository.Repository
}

type clients struct {
	account *client.AccountWrapper
}

type application struct {
	generator    *generator.Generator
	clients      clients
	connections  connections
	repositories repositories
	facade       *facade.Facade
	service      *token_service.Implementation

	cfg Config
	wg  *sync.WaitGroup
}

func newApp(ctx context.Context, cfg *Config) *application {
	var a = application{
		cfg: *cfg,
		wg:  &sync.WaitGroup{},
	}

	a.initClients(ctx)
	a.initGenerator()
	a.initConnections(ctx)
	a.initRepos()
	a.initFacade()
	a.initService()

	return &a
}

func (a *application) initClients(ctx context.Context) {
	conn, err := grpc.DialContext(ctx, a.cfg.Service.AccountAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("connect to account-service", zap.Error(err))
	}

	a.clients.account = client.NewAccountWrapper(account.NewAccountServiceClient(conn))
}

func (a *application) initGenerator() {
	a.generator = generator.NewGenerator(a.cfg.Service.TokenLength)
}

func (a *application) initConnections(ctx context.Context) {
	a.connections.memcached = memcache.NewCache[string, []byte]()

	if !a.cfg.Service.OnlyCacheMode {
		var err error
		a.connections.database, err = database.NewDatabase(ctx, a.cfg.Postgres.Dsn)
		if err != nil {
			logger.Error("init connect to database", zap.Error(err))
		}
	}
}

func (a *application) initRepos() {
	a.repositories = repositories{
		cache: cache.NewCache(a.connections.memcached),
		repo:  repository.New(a.connections.database),
	}
}

func (a *application) initFacade() {
	a.facade = facade.NewFacade(a.repositories.cache, a.repositories.repo, a.generator, a.cfg.Service.OnlyCacheMode)
}

func (a *application) initService() {
	a.service = token_service.NewService(a.facade, a.clients.account)
}

func (a *application) Run(ctx context.Context) error {
	server.StartTokenService(ctx, a.service, a.wg, a.cfg.Service.GrpcAddress)
	server.InitTokenSwagger(ctx, a.wg, a.cfg.Service.SwaggerAddress, a.cfg.Service.GrpcAddress)

	return nil
}

func (a *application) Close() {
	if !a.cfg.Service.OnlyCacheMode {
		a.connections.database.Close()
	}
}

func main() {
	configPath := flag.String("config", "", "config path")
	flag.Parse()

	config := readConfig(*configPath)

	ctx, cancel := context.WithCancel(context.Background())
	app := newApp(ctx, config)

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		sig := <-c
		logger.Info("received signal", zap.String("signal", sig.String()))
		cancel()
	}()

	defer app.Close()

	if err := app.Run(ctx); err != nil {
		logger.Fatal("can't run app", zap.Error(err))
	}

	app.wg.Wait()
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

	passwordSecret, err := yc.NewLockboxClient().WithTokenGetter(yc.ComputeMetadata()).GetSecret(cfg.PostgresPasswordLockboxSecretID)
	if err != nil {
		logger.Fatal("Failed to get secret from lockbox", zap.Error(err))
	}

	password := passwordSecret.Entries[0].TextValue
	cfg.Postgres.Password = password

	cfg.Postgres.Dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)

	return &cfg
}

type Config struct {
	Service                         TokenService `yaml:"service"`
	Postgres                        Postgres     `yaml:"postgres"`
	PostgresPasswordLockboxSecretID string       `yaml:"postgres_password_lockbox_secret_id"`
}

type TokenService struct {
	SwaggerAddress    string        `yaml:"swagger_address"`
	GrpcAddress       string        `yaml:"token_address"`
	AccountAddress    string        `yaml:"account_address"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
	OnlyCacheMode     bool          `yaml:"only_cache_mode"`
	TokenLength       int           `yaml:"token_length"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbname"`
	Password string
	Dsn      string
}
