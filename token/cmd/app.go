package main

import (
	"context"
	"sync"

	memcache "github.com/de1phin/iam/pkg/cache"

	token_service "github.com/de1phin/iam/token/api/token"
	"github.com/de1phin/iam/token/internal/cache"
	"github.com/de1phin/iam/token/internal/facade"
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
	connections  connections
	repositories repositories
	facade       *facade.Facade
	service      *token_service.Implementation

	wg *sync.WaitGroup
}

func newApp(ctx context.Context) *application {
	var a = application{}

	a.initConnections()
	a.initRepos()
	a.initFacade()
	a.initService()

	return &a
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
	a.facade = facade.NewFacade(a.repositories.cache, a.repositories.repo)
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
