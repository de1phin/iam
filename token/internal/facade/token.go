package facade

import (
	"context"

	"github.com/de1phin/iam/token/internal/model"
)

type Cache interface {
}

type Repository interface {
}

type Facade struct {
	cache Cache
	repo  Repository
}

func NewFacade(cache Cache, repo Repository) *Facade {
	return &Facade{
		cache: cache,
		repo:  repo,
	}
}

func (f *Facade) GenerateToken(ctx context.Context) (*model.Token, error) {
	return nil, nil
}

func (f *Facade) RefreshToken(ctx context.Context) (*model.Token, error) {
	return nil, nil
}

func (f *Facade) DeleteToken(ctx context.Context) error {
	return nil
}

func (f *Facade) CheckToken(context.Context) (bool, error) {
	return false, nil
}
