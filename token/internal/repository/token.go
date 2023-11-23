package repository

import (
	"context"

	"github.com/de1phin/iam/token/internal/model"
)

type Database interface {
}

type Repository struct {
	conn Database
}

func New(conn Database) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (c *Repository) GenerateToken(ctx context.Context) (*model.Token, error) {
	return nil, nil
}

func (c *Repository) RefreshToken(ctx context.Context) (*model.Token, error) {
	return nil, nil
}

func (c *Repository) DeleteToken(ctx context.Context) error {
	return nil
}

func (c *Repository) CheckToken(context.Context) (bool, error) {
	return false, nil
}
