package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/pkg/stringutil"
	"github.com/de1phin/iam/services/account/internal/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	generatedIdLength  = 10
	IdMinLength        = 4
	IdMaxLength        = 32
	NameMinLength      = 1
	NameMaxLength      = 64
	IdCharset          = "abcdefghijklmnopqrstuvwxyz"
	WellKnownIdCharset = "abcdefghijklmnopqrstuvwxyz._-"

	AccountNameFieldName = "name"
	AccountDescFieldName = "description"
	AccoundIdFieldName   = "id"
)

func generateNewId() string {
	id := ""
	for i := 0; i < generatedIdLength; i++ {
		id += string(IdCharset[rand.Intn(len(IdCharset))])
	}
	return id
}

func (s *AccountService) generateNewUniqueId(ctx context.Context) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			break
		}

		id := generateNewId()
		_, err := s.accounts.Get(id)
		if err == nil {
			continue
		}
		if errors.Is(err, database.ErrAlreadyExists{}) {
			return id, nil
		}
		return "", err
	}
}

func validateAccountFields(acc *account.Account) error {
	if len(acc.GetName()) < NameMinLength || len(acc.GetName()) > NameMaxLength {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("account name must be of length %d to %d", NameMinLength, NameMaxLength))
	}

	return nil
}

func (s *AccountService) validateNewAccount(acc *account.Account) error {
	for _, char := range acc.GetId() {
		if !bytes.ContainsRune([]byte(IdCharset), char) &&
			!bytes.ContainsRune([]byte(WellKnownIdCharset), char) {
			return status.Error(codes.InvalidArgument, "ID contains invalid symbol")
		}
	}
	_, err := s.accounts.Get(acc.GetId())
	if errors.Is(err, database.ErrAlreadyExists{}) {
		return status.Error(codes.AlreadyExists, "account with specified id already exists")
	}

	return validateAccountFields(acc)
}

func (s *AccountService) CreateAccount(ctx context.Context, req *account.CreateAccountRequest) (*account.CreateAccountResponse, error) {
	var err error
	id := req.GetWellKnownId()
	if id == "" {
		id, err = s.generateNewUniqueId(ctx)
		if err != nil {
			return nil, ErrorInternal()
		}
	}

	acc := &account.Account{
		Id:          id,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CreatedAt:   timestamppb.Now(),
	}

	err = s.validateNewAccount(acc)
	if err != nil {
		return nil, err
	}

	err = s.accounts.Create(acc)
	if err != nil {
		return nil, ErrorInternal()
	}
	return &account.CreateAccountResponse{
		Account: acc,
	}, nil
}

func (s *AccountService) GetAccount(_ context.Context, req *account.GetAccountRequest) (*account.GetAccountResponse, error) {
	acc, err := s.accounts.Get(req.GetAccountId())
	if err == nil {
		return &account.GetAccountResponse{
			Account: acc,
		}, nil
	}

	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}

	return nil, ErrorInternal()
}

func (s *AccountService) DeleteAccount(_ context.Context, req *account.DeleteAccountRequest) (*account.DeleteAccountResponse, error) {
	err := s.accounts.Delete(req.GetAccountId())
	if err == nil {
		return &account.DeleteAccountResponse{}, nil
	}

	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}

	return nil, ErrorInternal()
}

func updateAccountFields(acc *account.Account, req *account.UpdateAccountRequest) error {
	paths := req.GetUpdateMask().GetPaths()
	if stringutil.SliceContains(paths, AccoundIdFieldName) {
		return status.Error(codes.InvalidArgument, "ID cannot be updated")
	}
	if stringutil.SliceContains(paths, AccountNameFieldName) {
		acc.Name = req.GetName()
	}
	if stringutil.SliceContains(paths, AccountDescFieldName) {
		acc.Description = req.GetDescription()
	}
	return nil
}

func (s *AccountService) UpdateAccount(_ context.Context, req *account.UpdateAccountRequest) (*account.UpdateAccountResponse, error) {
	acc, err := s.accounts.Get(req.GetAccountId())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}
	if err != nil {
		return nil, ErrorInternal()
	}

	err = updateAccountFields(acc, req)
	if err != nil {
		return nil, err
	}

	err = validateAccountFields(acc)
	if err != nil {
		return nil, err
	}

	err = s.accounts.Update(acc)
	return &account.UpdateAccountResponse{
		Account: acc,
	}, nil
}
