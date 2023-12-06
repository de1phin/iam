package repository

import "time"

type Token struct {
	ID string `db:"id"`

	Ssh   string `db:"ssh"`
	Token string `db:"token"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}
