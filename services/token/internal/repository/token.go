package repository

import (
	"context"
	"errors"

	"github.com/de1phin/iam/pkg/database"
	"github.com/de1phin/iam/pkg/logger"
	"github.com/de1phin/iam/services/token/internal/model"
	"github.com/georgysavva/scany/pgxscan"
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
		CREATE TABLE IF NOT EXISTS token (
			token VARCHAR NOT NULL,
			account_id VARCHAR NOT NULL,
		
			expires_at TIMESTAMP NOT NULL
		);

		CREATE UNIQUE INDEX IF NOT EXISTS token_expires_at ON token(expires_at);
`

	_, err := db.Exec(ctx, query)
	return err
}

func (r *Repository) ExpireTokens(ctx context.Context) error {
	query := `DELETE FROM token WHERE expires_at < NOW() at time zone 'utc';`
	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *Repository) GetToken(ctx context.Context, token string) (*model.Token, error) {
	var tokenOrm Token

	query := `SELECT
			*
		FROM
			token
		WHERE
			token.token = $1`

	err := r.db.GetSingle(ctx, &tokenOrm, query, token)

	if pgxscan.NotFound(err) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &model.Token{
		Token:     tokenOrm.Token,
		AccountId: tokenOrm.AccountId,
		ExpiresAt: tokenOrm.ExpiresAt,
	}, err
}

func (r *Repository) CreateToken(ctx context.Context, token *model.Token) error {
	query := `
		INSERT INTO token
			(token, account_id, expires_at)
		VALUES
		    ($1, $2, $3);
`

	_, err := r.db.Exec(ctx, query, token.Token, token.AccountId, token.ExpiresAt)

	return err
}

func (r *Repository) DeleteToken(ctx context.Context, token string) error {
	query := `
		DELETE FROM token WHERE token.token = $1`

	_, err := r.db.Exec(ctx, query, token)

	if pgxscan.NotFound(err) {
		return ErrNotFound
	}

	return err
}

func (r *Repository) RefreshToken(ctx context.Context, token *model.Token) error {
	query := `
		UPDATE token SET expires_at = $1 WHERE token.token = $2`

	_, err := r.db.Exec(ctx, query, token.ExpiresAt, token.Token)
	if pgxscan.NotFound(err) {
		return ErrNotFound
	}
	return err
}
