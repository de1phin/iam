package database

import (
	account "github.com/de1phin/iam/genproto/services/account/api"
)

type AccountDatabase interface {
	Create(*account.Account) error
	Get(id string) (*account.Account, error)
	Update(*account.Account) error
	Delete(id string) error
}

type SshKeyDatabase interface {
	Create(*account.SshKey) error
	Get(fingerprint string) (*account.SshKey, error)
	List(accountId string) ([]*account.SshKey, error)
	Delete(accountId, fingerprint string) error
}

type ErrNotExist struct{}

func (ErrNotExist) Error() string {
	return "Object does not exist"
}

type ErrAlreadyExists struct{}

func (ErrAlreadyExists) Error() string {
	return "Object already exists"
}
