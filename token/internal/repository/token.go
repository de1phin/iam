package repository

import (
	"context"
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

func (r *Repository) GetToken(ctx context.Context, ssh string) (string, error) {
	return "", nil
}

func (r *Repository) SetToken(ctx context.Context, ssh string, token string) error {
	return nil
}

func (r *Repository) DeleteToken(ctx context.Context, ssh string) error {
	return nil
}

func (r *Repository) GetSsh(ctx context.Context, token string) (string, error) {
	return "", nil
}

func (r *Repository) SetSsh(ctx context.Context, token string, ssh string) error {
	return nil
}
