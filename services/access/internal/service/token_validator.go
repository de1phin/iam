package service

import "context"

//go:generate mockgen -source=token_validator.go -destination=token_validator_mock.go -package=service
type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (accountID string, isValid bool, err error)
}
