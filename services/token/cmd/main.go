package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	account "github.com/de1phin/iam/genproto/services/account/api"
	memcache "github.com/de1phin/iam/pkg/cache"
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

	wg *sync.WaitGroup
}

func newApp(ctx context.Context) *application {
	var a = application{
		wg: &sync.WaitGroup{},
	}

	a.initClients(ctx)
	a.initGenerator()
	a.initConnections()
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
	length := 512 // TODO config

	a.generator = generator.NewGenerator(length)
}

func (a *application) initConnections() {
	memcached := memcache.NewCache[string, []byte]()

	a.connections = connections{
		memcached: memcached,
	}
}

func (a *application) initRepos() {
	a.repositories = repositories{
		cache: cache.NewCache(a.connections.memcached),
	}
}

func (a *application) initFacade() {
	onlyCacheMod := true
	a.facade = facade.NewFacade(a.repositories.cache, a.repositories.repo, a.generator, onlyCacheMod)
}

func (a *application) initService() {
	a.service = token_service.NewService(a.facade, a.clients.account)
}

func (a *application) Run(ctx context.Context) error {
	host := "token-service" // TODO
	server.StartTokenService(ctx, a.service, a.wg, host)

	return nil
}

func (a *application) Close() {
	a.wg.Wait()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app := newApp(ctx)

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
