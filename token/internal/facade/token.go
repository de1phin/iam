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
	DeleteToken(ctx context.Context, ssh string) error

	GetExist(ctx context.Context, token string) (bool, error)
	SetExist(ctx context.Context, token string, isExist bool) error
}

type Repository interface {
	GetToken(ctx context.Context, ssh string) (string, error)
	SetToken(ctx context.Context, ssh string, token string) error
	DeleteToken(ctx context.Context, ssh string) error

	GetExist(ctx context.Context, token string) (bool, error)
}

type Generator interface {
	Generate() string
}

type Facade struct {
	cache     Cache
	repo      Repository
	generator Generator
	onlyCache bool
}

func NewFacade(cache Cache, repo Repository, generator Generator, onlyCacheMod bool) *Facade {
	return &Facade{
		cache:     cache,
		repo:      repo,
		generator: generator,
		onlyCache: onlyCacheMod,
	}
}

func (f *Facade) GenerateToken(ctx context.Context, ssh string) (*model.Token, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/GenerateToken")
	defer span.Finish()

	if f.onlyCache {
		var isNotFound bool

		token, err := f.cache.GetToken(ctx, ssh)
		if err != nil {
			isNotFound = errors.Is(err, cache.ErrNotFound)
			if !isNotFound {
				return nil, err
			}
		}

		if !isNotFound {
			isExist, err := f.cache.GetExist(ctx, token)
			if err != nil {
				return nil, err
			}

			if isExist {
				return convertToModelToken(token), err
			}
		}

		token = f.generator.Generate()

		if err = f.cache.SetToken(ctx, ssh, token); err != nil {
			return nil, err
		}
		if err = f.cache.SetExist(ctx, token, true); err != nil {
			return nil, err
		}
		return convertToModelToken(token), nil
	}

	// TODO добавить логику с базой

	return nil, nil
}

func (f *Facade) RefreshToken(ctx context.Context, ssh string) (*model.Token, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/RefreshToken")
	defer span.Finish()

	token := f.generator.Generate()

	if f.onlyCache {
		if err := f.cache.SetToken(ctx, ssh, token); err != nil {
			return nil, err
		}
		if err := f.cache.SetExist(ctx, token, true); err != nil {
			return nil, err
		}
		return convertToModelToken(token), nil
	}

	if err := f.cache.SetToken(ctx, ssh, token); err != nil {
		logger.Error("cache SetToken", zap.Error(err))
	}

	if err := f.cache.SetExist(ctx, token, true); err != nil {
		logger.Error("cache SetToken", zap.Error(err))
	}

	if err := f.repo.SetToken(ctx, ssh, token); err != nil {
		return nil, err
	}

	return convertToModelToken(token), nil
}

func (f *Facade) DeleteToken(ctx context.Context, ssh string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/DeleteToken")
	defer span.Finish()

	if f.onlyCache {
		token, err := f.cache.GetToken(ctx, ssh)
		if err != nil {
			return err
		}

		if err = f.cache.DeleteToken(ctx, ssh); err != nil {
			return err
		}

		return f.cache.SetExist(ctx, token, false)
	}

	if err := f.cache.DeleteToken(ctx, ssh); err != nil {
		logger.Error("cache DeleteToken", zap.Error(err))
	}

	return f.repo.DeleteToken(ctx, ssh)
}

func (f *Facade) CheckToken(ctx context.Context, tk model.Token) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/CheckToken")
	defer span.Finish()

	token := tk.Token

	if f.onlyCache {
		isExist, err := f.cache.GetExist(ctx, token)
		if err != nil {
			if errors.Is(err, cache.ErrNotFound) {
				return false, nil
			}
		}
		return isExist, err
	}

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
