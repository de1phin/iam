package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	account "github.com/de1phin/iam/genproto/services/account/api"
	memcache "github.com/de1phin/iam/pkg/cache"
	"github.com/de1phin/iam/pkg/database"
	"github.com/de1phin/iam/pkg/logger"
	token_service "github.com/de1phin/iam/services/token/app/token"
	"github.com/de1phin/iam/services/token/internal/cache"
	"github.com/de1phin/iam/services/token/internal/client"
	"github.com/de1phin/iam/services/token/internal/facade"
	"github.com/de1phin/iam/services/token/internal/generator"
	"github.com/de1phin/iam/services/token/internal/repository"
	"github.com/de1phin/iam/services/token/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
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
	conn, err := grpc.DialContext(ctx, a.cfg.Service.AccountAddress)
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
		a.connections.database, err = database.NewDatabase(ctx, a.cfg.Service.Dsn)
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
	server.StartTokenService(ctx, a.service, a.wg, a.cfg.Service.TokenAddress)

	return nil
}

func (a *application) Close() {
	a.connections.database.Close()

	a.wg.Wait()
}

func main() {
	configPath := flag.String("config", "", "config path")
	flag.Parse()

	config := readConfig(*configPath)

	ctx, cancel := context.WithCancel(context.Background())
	app := newApp(ctx, config)

	app.wg.Add(1)
	go func() {
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

type Config struct {
	Service TokenService `yaml:"service"`
}

type TokenService struct {
	TokenAddress      string        `yaml:"token_address"`
	AccountAddress    string        `yaml:"account_address"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout"`
	Dsn               string        `yaml:"token_dsn"`
	OnlyCacheMode     bool          `yaml:"only_cache_mode"`
	TokenLength       int           `yaml:"token_length"`
}
