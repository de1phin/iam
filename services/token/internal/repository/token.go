package repository

import (
	"context"
	"errors"

	"github.com/de1phin/iam/pkg/database"
	"github.com/de1phin/iam/pkg/logger"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db database.Database
}

func New(ctx context.Context, db database.Database) *Repository {
	if err := initSchema(ctx, db); err != nil {
		logger.Error("init token database", zap.Error(err))
	}
	return &Repository{
		db: db,
	}
}

func initSchema(ctx context.Context, db database.Database) error {
	query := `
		CREATE TABLE token_info (
			id SERIAL PRIMARY KEY,
		
			ssh TEXT NOT NULL,
			token TEXT NOT NULL,
		
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now(),
			deleted_at TIMESTAMP
		);
		
		CREATE UNIQUE INDEX token_info_token ON token_info(token);
`

	_, err := db.Exec(ctx, query)
	return err
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
