package repository

import (
	"context"
	"errors"

	"github.com/de1phin/iam/pkg/database"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db database.Database
}

func New(db database.Database) *Repository {
	return &Repository{
		db: db,
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

	err := r.db.GetSingle(ctx, &token, query, ssh)

	return token.Token, err
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
		    updated_at = NOW() at time zone 'utc',
		    deleted_at = NULL;
`

	_, err := r.db.Exec(ctx, query, token, ssh)

	return err
}

func (r *Repository) DeleteToken(ctx context.Context, ssh string) error {
	query := `
		UPDATE token_info 
		SET deleted_at = NOW() at time zone 'utc' 
		WHERE ssh = $1;`

	_, err := r.db.Exec(ctx, query, ssh)

	return err
}

func (r *Repository) GetSsh(ctx context.Context, token string) (string, error) {
	var repoToken Token

	query := `SELECT
			ssh
		FROM
			token_info
		WHERE
			token = $1`

	err := r.db.GetSingle(ctx, &repoToken, query, token)

	return repoToken.Ssh, err
}
