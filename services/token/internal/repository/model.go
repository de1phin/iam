package repository

import "time"

type Token struct {
	Token     string `db:"token"`
	AccountId string `db:"account_id"`

	ExpiresAt time.Time `db:"expires_at"`
}
