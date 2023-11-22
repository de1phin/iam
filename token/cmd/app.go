package main

import (
	"context"
	"sync"

	token_service "github.com/de1phin/iam/token/api/token"
	"github.com/de1phin/iam/token/internal/server"
)

type connections struct {
}

type repositories struct {
}

type providers struct {
}

type services struct {
	token *token_service.Implementation
}

type application struct {
	connections  connections
	repositories repositories
	providers    providers
	services     services

	wg *sync.WaitGroup
}

func newApp(ctx context.Context) *application {
	var a = application{}

	a.initConnections()
	a.initRepos()
	a.initProviders()
	a.initService()

	return &a
}

func (a *application) initConnections() {
	a.connections = connections{}
}

func (a *application) initProviders() {
	a.providers = providers{}
}

func (a *application) initService() {
	a.services = services{
		token: token_service.NewService(),
	}
}

func (a *application) initRepos() {
	a.repositories = repositories{}
}

func (a *application) Run(ctx context.Context) error {
	server.StartTokenService(ctx, a.services.token, a.wg)

	return nil
}

func (a *application) Close() {
	a.wg.Wait()
}
