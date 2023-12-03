//go:generate mockery --case=underscore --name Cache --with-expecter=true
//go:generate mockery --case=underscore --name Repository --with-expecter=true
//go:generate mockery --case=underscore --name Generator --with-expecter=true

package facade

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"

	"github.com/de1phin/iam/token/internal/cache"
	"github.com/de1phin/iam/token/internal/model"
)

type Cache interface {
	GetToken(ctx context.Context, ssh string) (string, error)
	SetToken(ctx context.Context, ssh string, token string) error
	DeleteToken(ctx context.Context, ssh string) error

	GetSsh(ctx context.Context, token string) (string, error)
	SetSsh(ctx context.Context, token string, ssh string) error
	DeleteSsh(ctx context.Context, ssh string) error
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
		var isNotFound bool

		token, err := f.cache.GetToken(ctx, ssh)
		if err != nil {
			isNotFound = errors.Is(err, cache.ErrNotFound)
			if !isNotFound {
				return nil, err
			}
		}

		if !isNotFound {
			cacheSsh, err := f.cache.GetSsh(ctx, token)
			if err != nil {
				return nil, err
			}

			if cacheSsh != "" {
				return convertToModelToken(token), err
			}
		}

		token = f.generator.Generate()

		if err = f.cache.SetToken(ctx, ssh, token); err != nil {
			return nil, err
		}
		if err = f.cache.SetSsh(ctx, token, ssh); err != nil {
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
		if err := f.cache.SetSsh(ctx, token, ssh); err != nil {
			return nil, err
		}
		return convertToModelToken(token), nil
	}

	//if err := f.cache.SetToken(ctx, ssh, token); err != nil {
	//	logger.Error("cache SetToken", zap.Error(err))
	//}
	//
	//if err := f.cache.SetSsh(ctx, token, ssh); err != nil {
	//	logger.Error("cache SetToken", zap.Error(err))
	//}
	//
	//if err := f.repo.SetToken(ctx, ssh, token); err != nil {
	//	return nil, err
	//}
	//
	//return convertToModelToken(token), nil

	return nil, nil
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

		return f.cache.DeleteSsh(ctx, token)
	}

	//if err := f.cache.DeleteToken(ctx, ssh); err != nil {
	//	logger.Error("cache DeleteToken", zap.Error(err))
	//}
	//
	//return f.repo.DeleteToken(ctx, ssh)

	return nil
}

func (f *Facade) GetSshByToken(ctx context.Context, tk model.Token) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/GetSshByToken")
	defer span.Finish()

	token := tk.Token

	if f.onlyCache {
		return f.cache.GetSsh(ctx, token)
	}

	//ssh, err := f.cache.GetSsh(ctx, token)
	//if err == nil {
	//	return ssh, nil
	//}
	//if ssh == "" {
	//	logger.Error("cache GetExist", zap.Error(err))
	//}
	//
	//ssh, err = f.repo.GetSsh(ctx, token)
	//if err != nil {
	//	return "", fmt.Errorf("repo GetExist: %w", err)
	//}
	//
	//if err = f.cache.SetSsh(ctx, token, ssh); err != nil {
	//	logger.Error("cache SetExist", zap.Error(err))
	//}
	//
	//return ssh, nil

	return "", nil
}

func convertToModelToken(token string) *model.Token {
	return &model.Token{Token: token}
}
