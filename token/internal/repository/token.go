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

func (r *Repository) GenerateToken(ctx context.Context) (*model.Token, error) {
	return nil, nil
}

func (r *Repository) GetToken(ctx context.Context, ssh string) (string, error) {
	return "", nil
}

func (r *Repository) SetToken(ctx context.Context, ssh string, token string) error {
	return nil
}

func (r *Repository) GetExist(ctx context.Context, token string) (bool, error) {
	return false, nil
}

func (r *Repository) SetExist(ctx context.Context, token string, isExist bool) error {
	return nil
}
