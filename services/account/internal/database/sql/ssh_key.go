package sql

import (
	"context"
	"time"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SqlSshKeyDatabase struct {
	db *SqlDatabase
}

type SshKey struct {
	Fingerprint string    `db:"fingerprint"`
	AccountId   string    `db:"account_id"`
	PublicKey   string    `db:"public_key"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewSqlSshKeyDatabase(db *SqlDatabase) *SqlSshKeyDatabase {
	return &SqlSshKeyDatabase{
		db: db,
	}
}

func (r *SqlSshKeyDatabase) Get(fingerprint string) (*account.SshKey, error) {
	var res SshKey

	query := `
		SELECT
			*
		FROM
			ssh_key
		WHERE
			fingerprint = $1`

	err := r.db.GetSingle(context.Background(), &res, query, fingerprint)

	if err != nil {
		return nil, convertError(err)
	}

	return &account.SshKey{
		Fingerprint: res.Fingerprint,
		AccountId:   res.AccountId,
		PublicKey:   res.PublicKey,
		CreatedAt:   timestamppb.New(res.CreatedAt),
	}, nil
}

func (r *SqlSshKeyDatabase) List(accountId string) ([]*account.SshKey, error) {
	var res []SshKey

	query := `
		SELECT
			*
		FROM
			ssh_key
		WHERE account_id = $1
	`

	_, err := r.db.Exec(context.Background(), query, res)
	if err != nil {
		return nil, convertError(err)
	}

	keys := make([]*account.SshKey, len(res))
	for i, key := range res {
		keys[i] = &account.SshKey{
			Fingerprint: key.Fingerprint,
			AccountId:   key.AccountId,
			PublicKey:   key.PublicKey,
			CreatedAt:   timestamppb.New(key.CreatedAt),
		}
	}
	return keys, nil
}

func (r *SqlSshKeyDatabase) Create(key *account.SshKey) error {
	query := `
		INSERT INTO ssh_key(fingerprint, account_id, public_key, created_at)
		VALUES($1, $2, $3, $4)
	`

	_, err := r.db.Exec(context.Background(), query,
		key.Fingerprint, key.AccountId, key.PublicKey, key.CreatedAt.AsTime())
	return convertError(err)
}

func (r *SqlSshKeyDatabase) Delete(accountId string, fingerprint string) error {
	query := `
		DELETE FROM ssh_key WHERE account_id = $1 AND fingerprint = $2
	`

	_, err := r.db.Exec(context.Background(), query, accountId, fingerprint)
	return convertError(err)
}
