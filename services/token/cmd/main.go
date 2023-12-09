package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	wg           *sync.WaitGroup
	onlyCacheMod bool
}

func newApp(ctx context.Context) *application {
	var a = application{
		wg: &sync.WaitGroup{},
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
	host := "account-service" // TODO
	conn, err := grpc.DialContext(ctx, host)
	if err != nil {
		logger.Fatal("connect to account-service", zap.Error(err))
	}

	a.clients.account = client.NewAccountWrapper(account.NewAccountServiceClient(conn))
}

func (a *application) initGenerator() {
	length := 100 // TODO config

	a.generator = generator.NewGenerator(length)
}

func (a *application) initConnections(ctx context.Context) {
	memcached := memcache.NewCache[string, []byte]()

	dsn := "" // TODO
	db, err := database.NewDatabase(ctx, dsn)
	if err != nil {
		logger.Error("init connect to database, only cached mode turn on", zap.Error(err))
		a.onlyCacheMod = true
	}

	a.connections = connections{
		memcached: memcached,
		database:  db,
	}
}

func (a *application) initRepos() {
	a.repositories = repositories{
		cache: cache.NewCache(a.connections.memcached),
		repo:  repository.New(a.connections.database),
	}
}

func (a *application) initFacade() {
	a.onlyCacheMod = true
	a.facade = facade.NewFacade(a.repositories.cache, a.repositories.repo, a.generator, a.onlyCacheMod)
}

func (a *application) initService() {
	a.service = token_service.NewService(a.facade, a.clients.account)
}

func (a *application) Run(ctx context.Context) error {
	host := ":8080" // TODO
	server.StartTokenService(ctx, a.service, a.wg, host)
	server.InitTokenSwagger(ctx, a.wg, host)

	return nil
}

func (a *application) Close() {
	a.connections.database.Close()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app := newApp(ctx)

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
