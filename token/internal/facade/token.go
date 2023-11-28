package facade

import (
	"context"
	"errors"
	"fmt"

	"github.com/de1phin/iam/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/de1phin/iam/token/internal/cache"
	"github.com/de1phin/iam/token/internal/model"
)

type Cache interface {
	GetToken(ctx context.Context, ssh string) (string, error)
	SetToken(ctx context.Context, ssh string, token string) error

	GetExist(ctx context.Context, token string) (bool, error)
	SetExist(ctx context.Context, token string, isExist bool) error
}

type Repository interface {
	GetToken(ctx context.Context, ssh string) (string, error)
	SetToken(ctx context.Context, ssh string, token string) error

	GetExist(ctx context.Context, token string) (bool, error)
	SetExist(ctx context.Context, token string, isExist bool) error
}

type Generator interface {
	Generate() string
}

type Facade struct {
	cache     Cache
	repo      Repository
	generator Generator
}

func NewFacade(cache Cache, repo Repository, generator Generator) *Facade {
	return &Facade{
		cache:     cache,
		repo:      repo,
		generator: generator,
	}
}

func (f *Facade) GenerateToken(ctx context.Context, ssh string) (*model.Token, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/GenerateToken")
	defer span.Finish()

	token, err := f.cache.GetToken(ctx, ssh)
	if err == nil {
		return convertToModelToken(token), nil
	}

	return nil, err
}

func (f *Facade) RefreshToken(ctx context.Context, ssh string) (*model.Token, error) {
	return nil, nil
}

func (f *Facade) DeleteToken(ctx context.Context, ssh string) error {
	return nil
}

func (f *Facade) CheckToken(ctx context.Context, tk model.Token) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/GenerateToken")
	defer span.Finish()

	token := tk.Token

	isExist, err := f.cache.GetExist(ctx, token)
	if err == nil {
		return isExist, nil
	}
	if !errors.Is(err, cache.ErrNotFound) {
		logger.Error("cache GetExist", zap.Error(err))
	}

	isExist, err = f.repo.GetExist(ctx, token)
	if err != nil {
		return false, fmt.Errorf("repo GetExist: %w", err)
	}

	if err = f.cache.SetExist(ctx, token, isExist); err != nil {
		logger.Error("cache SetExist", zap.Error(err))
	}

	return isExist, nil
}

func convertToModelToken(token string) *model.Token {
	return &model.Token{Token: token}
}
