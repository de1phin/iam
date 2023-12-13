package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/de1phin/iam/pkg/database"
	accountServiceDatabase "github.com/de1phin/iam/services/account/internal/database"
	"github.com/georgysavva/scany/pgxscan"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SqlDatabase struct {
	database.Database
}

func timeToStr(t *timestamppb.Timestamp) string {
	return t.AsTime().Format(time.RFC1123)
}

func strToTime(s string) *timestamppb.Timestamp {
	t, _ := time.Parse(time.RFC1123, s)
	return timestamppb.New(t)
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
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
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
		
			created_at TIMESTAMP NOT NULL DEFAULT now()
		);
		
		CREATE TABLE IF NOT EXISTS ssh_key (
			fingerprint VARCHAR UNIQUE NOT NULL PRIMARY KEY,
			public_key VARCHAR NOT NULL,
			account_id VARCHAR NOT NULL,
		
			created_at TIMESTAMP NOT NULL DEFAULT now()
		);

		CREATE UNIQUE INDEX IF NOT EXISTS ssh_key_account_id ON ssh_key(account_id);
`

	_, err := db.Exec(ctx, query)
	return err
}
