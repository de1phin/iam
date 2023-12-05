//go:generate mockery --case=underscore --name Cache --with-expecter=true
//go:generate mockery --case=underscore --name Repository --with-expecter=true
//go:generate mockery --case=underscore --name Generator --with-expecter=true

package facade

import (
	"context"
	"errors"
	"fmt"

	"github.com/de1phin/iam/pkg/logger"
	"github.com/de1phin/iam/services/token/internal/cache"
	"github.com/de1phin/iam/services/token/internal/model"
	"github.com/de1phin/iam/services/token/internal/repository"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Cache interface {
	GetToken(ctx context.Context, ssh string) (string, error)
	SetToken(ctx context.Context, ssh string, token string) error
	DeleteToken(ctx context.Context, ssh string) error

	GetSsh(ctx context.Context, token string) (string, error)
	SetSsh(ctx context.Context, token string, ssh string) error
	DeleteSsh(ctx context.Context, token string) error
}

type Repository interface {
	GetToken(ctx context.Context, ssh string) (string, error)
	SetToken(ctx context.Context, ssh string, token string) error
	DeleteToken(ctx context.Context, ssh string) error

	GetSsh(ctx context.Context, token string) (string, error)
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
		var tokenNotFound bool

		token, err := f.cache.GetToken(ctx, ssh)
		if err != nil {
			tokenNotFound = errors.Is(err, cache.ErrNotFound)
			if !tokenNotFound {
				return nil, err
			}
		}

		if !tokenNotFound {
			_, err = f.cache.GetSsh(ctx, token)
			if err == nil {
				return convertToModelToken(token), err
			}

			if !errors.Is(err, cache.ErrNotFound) {
				return nil, err
			}
		}

		newToken := f.generator.Generate()

		if err = f.cache.SetToken(ctx, ssh, newToken); err != nil {
			return nil, err
		}
		if err = f.cache.DeleteSsh(ctx, token); err != nil {
			return nil, err
		}
		if err = f.cache.SetSsh(ctx, newToken, ssh); err != nil {
			return nil, err
		}

		return convertToModelToken(newToken), nil
	}

	var isNotFound bool

	token, err := f.cache.GetToken(ctx, ssh)
	if err != nil {
		isNotFound = errors.Is(err, cache.ErrNotFound)
		if !isNotFound {
			logger.Error("cache GetToken", zap.Error(err))
		}
	}

	if !isNotFound && err == nil {
		_, err = f.cache.GetSsh(ctx, token)
		if err == nil {
			return convertToModelToken(token), err
		}
		if !errors.Is(err, cache.ErrNotFound) {
			logger.Error("cache GetSsh", zap.Error(err))
		}
	}

	token, err = f.repo.GetToken(ctx, ssh)
	if err == nil {
		return convertToModelToken(token), nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	newToken := f.generator.Generate()

	if err = f.cache.SetToken(ctx, ssh, newToken); err != nil {
		logger.Error("cache SetToken", zap.Error(err))
	}
	if err = f.cache.DeleteSsh(ctx, token); err != nil {
		logger.Error("cache DeleteSsh", zap.Error(err))
	}
	if err = f.cache.SetSsh(ctx, newToken, ssh); err != nil {
		logger.Error("cache SetSsh", zap.Error(err))
	}

	if err = f.repo.SetToken(ctx, ssh, newToken); err != nil {
		return nil, err
	}

	return convertToModelToken(newToken), nil
}

func (f *Facade) RefreshToken(ctx context.Context, ssh string) (*model.Token, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/RefreshToken")
	defer span.Finish()

	token := f.generator.Generate()

	if f.onlyCache {
		if err := f.cache.SetToken(ctx, ssh, token); err != nil {
			return nil, err
		}
		if err := f.cache.SetSsh(ctx, token, ssh); err != nil {
			return nil, err
		}
		return convertToModelToken(token), nil
	}

	if err := f.cache.SetToken(ctx, ssh, token); err != nil {
		logger.Error("cache SetToken", zap.Error(err))
	}

	oldToken, err := f.cache.GetToken(ctx, ssh)
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			logger.Error("cache GetToken", zap.Error(err))
		}
	} else {
		if err = f.cache.DeleteSsh(ctx, oldToken); err != nil {
			logger.Error("cache DeleteSsh", zap.Error(err))
		}
	}

	if err = f.cache.SetSsh(ctx, token, ssh); err != nil {
		logger.Error("cache SetSsh", zap.Error(err))
	}

	if err = f.repo.SetToken(ctx, ssh, token); err != nil {
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
			if errors.Is(err, cache.ErrNotFound) {
				return nil
			}
			return err
		}

		if err = f.cache.DeleteToken(ctx, ssh); err != nil {
			return err
		}

		return f.cache.DeleteSsh(ctx, token)
	}

	token, err := f.cache.GetToken(ctx, ssh)
	if err == nil {
		if err = f.cache.DeleteSsh(ctx, token); err != nil {
			logger.Error("cache DeleteToken", zap.Error(err))
		}
	} else {
		if !errors.Is(err, cache.ErrNotFound) {
			logger.Error("cache GetToken", zap.Error(err))
		}
	}

	if err = f.cache.DeleteToken(ctx, ssh); err != nil {
		logger.Error("cache DeleteToken", zap.Error(err))
	}

	return f.repo.DeleteToken(ctx, ssh)
}

func (f *Facade) GetSshByToken(ctx context.Context, tk model.Token) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/GetSshByToken")
	defer span.Finish()

	token := tk.Token

	if f.onlyCache {
		return f.cache.GetSsh(ctx, token)
	}

	ssh, err := f.cache.GetSsh(ctx, token)
	if err == nil {
		return ssh, nil
	}
	notFound := errors.Is(err, cache.ErrNotFound)
	if !notFound {
		logger.Error("cache GetSsh", zap.Error(err))
	}

	ssh, err = f.repo.GetSsh(ctx, token)
	if err != nil {
		return "", fmt.Errorf("repo GetSsh: %w", err)
	}

	if notFound {
		if err = f.cache.SetSsh(ctx, token, ssh); err != nil {
			logger.Error("cache SetSsh", zap.Error(err))
		}
		if err = f.cache.SetToken(ctx, ssh, token); err != nil {
			logger.Error("cache SetSsh", zap.Error(err))
		}
	}

	return ssh, nil
}

func convertToModelToken(token string) *model.Token {
	return &model.Token{Token: token}
}
