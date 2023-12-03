package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	memcache "github.com/de1phin/iam/pkg/cache"
	"github.com/de1phin/iam/pkg/logger"
	"go.uber.org/zap"

	token_service "github.com/de1phin/iam/token/api/token"
	"github.com/de1phin/iam/token/internal/cache"
	"github.com/de1phin/iam/token/internal/facade"
	"github.com/de1phin/iam/token/internal/generator"
	"github.com/de1phin/iam/token/internal/repository"
	"github.com/de1phin/iam/token/internal/server"
)

type connections struct {
	memcached *memcache.Cache[string, []byte]
}

type repositories struct {
	cache *cache.MemCache
	repo  *repository.Repository
}

type application struct {
	generator    *generator.Generator
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

	a.initGenerator()
	a.initConnections()
	a.initRepos()
	a.initFacade()
	a.initService()

	return &a
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
	a.service = token_service.NewService(a.facade)
}

func (a *application) Run(ctx context.Context) error {
	server.StartTokenService(ctx, a.service, a.wg)

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
