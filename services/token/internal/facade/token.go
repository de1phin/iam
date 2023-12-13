//go:generate mockery --case=underscore --name Cache --with-expecter=true
//go:generate mockery --case=underscore --name Repository --with-expecter=true
//go:generate mockery --case=underscore --name Generator --with-expecter=true

package facade

import (
	"context"
	"time"

	"github.com/de1phin/iam/pkg/logger"
	"github.com/de1phin/iam/services/token/internal/model"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type Repository interface {
	GetToken(ctx context.Context, token string) (*model.Token, error)
	CreateToken(ctx context.Context, token *model.Token) error
	ExpireTokens(ctx context.Context) error
	DeleteToken(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, token *model.Token) error
}

type Generator interface {
	Generate() string
}

type Facade struct {
	repo      Repository
	generator Generator
}

const TokenDefaultValidity = time.Hour * 12

func NewFacade(repo Repository, generator Generator) *Facade {
	return &Facade{
		repo:      repo,
		generator: generator,
	}
}

func (f *Facade) ExpireTokens(ctx context.Context) error {
	return f.repo.ExpireTokens(ctx)
}

func (f *Facade) CreateToken(ctx context.Context, accountId string) (*model.Token, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/CreateToken")
	defer span.Finish()

	token := f.generator.Generate()

	modelToken := &model.Token{
		Token:     token,
		AccountId: accountId,
		ExpiresAt: time.Now().Add(TokenDefaultValidity).UTC(),
	}

	err := f.repo.CreateToken(ctx, modelToken)
	if err != nil {
		logger.Error("repo.SetToken", zap.Error(err))
		return nil, err
	}

	return modelToken, nil
}

func (f *Facade) RefreshToken(ctx context.Context, token string) (*model.Token, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/RefreshToken")
	defer span.Finish()

	modelToken, err := f.repo.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}

	modelToken.ExpiresAt = time.Now().Add(TokenDefaultValidity).UTC()

	err = f.repo.RefreshToken(ctx, modelToken)

	if err != nil {
		return nil, err
	}

	return modelToken, err
}

func (f *Facade) DeleteToken(ctx context.Context, token string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/DeleteToken")
	defer span.Finish()

	return f.repo.DeleteToken(ctx, token)
}

func (f *Facade) GetToken(ctx context.Context, token string) (*model.Token, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "facade/ExchangeToken")
	defer span.Finish()

	return f.repo.GetToken(ctx, token)
}
