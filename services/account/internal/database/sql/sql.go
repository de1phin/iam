package sql

import (
	"context"

	"github.com/de1phin/iam/pkg/database"
	accountServiceDatabase "github.com/de1phin/iam/services/account/internal/database"
	"github.com/georgysavva/scany/pgxscan"
)

type SqlDatabase struct {
	database.Database
}

func convertError(err error) error {
	if pgxscan.NotFound(err) {
		return accountServiceDatabase.ErrNotExist{}
	}

	return err
}

func New(ctx context.Context, db database.Database) (*SqlDatabase, error) {
	err := initSchema(ctx, db)
	if err != nil {
		return nil, err
	}
	return &SqlDatabase{
		Database: db,
	}, nil
}

func initSchema(ctx context.Context, db database.Database) error {
	query := `
		CREATE TABLE IF NOT EXISTS account (
			id VARCHAR UNIQUE NOT NULL PRIMARY KEY,

			name VARCHAR NOT NULL,
			description VARCHAR,
		
			created_at TIMESTAMP NOT NULL DEFAULT now(),
		);
		
		CREATE TABLE IF NOT EXISTS ssh_key (
			fingerprint VARCHAR UNIQUE NOT NULL PRIMARY KEY,
			public_key VARCHAR UNIQUE NOT NULL,
			account_id VARCHAR NOT NULL,
		
			created_at TIMESTAMP NOT NULL DEFAULT now(),
		);

		CREATE UNIQUE INDEX IF NOT EXISTS ssh_key_account_id ON ssh_key(account_id);
`

	_, err := db.Exec(ctx, query)
	return err
}
