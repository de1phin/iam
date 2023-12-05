package repository

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

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
	var token Token

	query := `SELECT
			token
		FROM
			token_info
		WHERE
			ssh = $1`

	query = query

	return token.Token, nil
}

func (r *Repository) SetToken(ctx context.Context, ssh string, token string) error {
	query := `
		INSERT INTO token_info
			(token, ssh)
		VALUES
		    ($1, $2)
		ON CONFLICT (token) DO UPDATE SET 
		    ssh = EXCLUDED.ssh,
		    token = EXCLUDED.token,
		    updated_at = NOW() at time zone 'utc';
`

	query = query

	return nil
}

func (r *Repository) DeleteToken(ctx context.Context, ssh string) error {
	query := `
		UPDATE token_info 
		SET deleted_at = NOW() at time zone 'utc' 
		WHERE ssh = $1;`

	query = query

	return nil
}

func (r *Repository) GetSsh(ctx context.Context, token string) (string, error) {
	var repoToken Token

	query := `SELECT
			ssh
		FROM
			token_info
		WHERE
			token = $1`

	query = query

	return repoToken.Ssh, nil
}
