package sql

import (
	"context"
	"time"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SqlAccountDatabase struct {
	db *SqlDatabase
}

type Account struct {
	Id          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewSqlAccountDatabase(db *SqlDatabase) *SqlAccountDatabase {
	return &SqlAccountDatabase{
		db: db,
	}
}

func (r *SqlAccountDatabase) Get(id string) (*account.Account, error) {
	var res Account

	query := `
		SELECT
			*
		FROM
			account
		WHERE
			id = $1`

	err := r.db.GetSingle(context.Background(), &res, query, id)

	if err != nil {
		return nil, convertError(err)
	}

	return &account.Account{
		Id:          res.Id,
		Name:        res.Name,
		Description: res.Description,
		CreatedAt:   timestamppb.New(res.CreatedAt),
	}, nil
}

func (r *SqlAccountDatabase) Create(account *account.Account) error {
	query := `
		INSERT INTO account(id, name, description, created_at)
		VALUES($1, $2, $3, $4)
	`

	_, err := r.db.Exec(context.Background(), query,
		account.Id, account.Name, account.Description, account.CreatedAt.AsTime())
	return convertError(err)
}

func (r *SqlAccountDatabase) Update(account *account.Account) error {
	query := `
		UPDATE account
		SET name = $1, description = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(context.Background(), query,
		account.Name, account.Description, account.Id)
	return convertError(err)
}

func (r *SqlAccountDatabase) Delete(id string) error {
	query := `
		DELETE FROM account WHERE id = $1
	`

	_, err := r.db.Exec(context.Background(), query, id)
	return convertError(err)
}
