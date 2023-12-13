package model

import "time"

type Token struct {
	Token     string
	AccountId string
	ExpiresAt time.Time
}
