package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/pkg/sshutil"
	"github.com/de1phin/iam/pkg/stringutil"
	"github.com/de1phin/iam/services/account/internal/database"
	"github.com/de1phin/iam/services/account/pkg/ctxlog"
	"go.uber.org/zap"
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
			ctxlog.Logger(ctx).Warn("generateNewUniqueId canceled")
			return "", ctx.Err()
		default:
			break
		}

		id := generateNewId()
		_, err := s.accounts.Get(id)
		if err == nil {
			continue
		}
		if errors.Is(err, database.ErrNotExist{}) {
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

func (s *AccountService) validateNewAccount(ctx context.Context, acc *account.Account) error {
	for _, char := range acc.GetId() {
		if !bytes.ContainsRune([]byte(IdCharset), char) &&
			!bytes.ContainsRune([]byte(WellKnownIdCharset), char) {
			return status.Error(codes.InvalidArgument, "ID contains invalid symbol")
		}
	}

	_, err := s.accounts.Get(acc.GetId())
	if err == nil {
		return status.Error(codes.AlreadyExists, "account with specified id already exists")
	}
	if err != nil && !errors.Is(err, database.ErrNotExist{}) {
		ctxlog.Logger(ctx).With(zap.String("account_id", acc.GetId())).Error("Internal error on AccountDatabase.Get", zap.Error(err))
		return ErrorInternal(err)
	}

	return validateAccountFields(acc)
}

func (s *AccountService) CreateAccount(ctx context.Context, req *account.CreateAccountRequest) (*account.CreateAccountResponse, error) {
	logger := s.logger.Named("CreateAccount")

	var err error
	id := req.GetWellKnownId()
	if id == "" {
		id, err = s.generateNewUniqueId(ctxlog.ContextWithLogger(ctx, logger))
		if err != nil {
			logger.Error("generateNewUniqueId failed", zap.Error(err))
			return nil, ErrorInternal(err)
		}
	}
	logger = logger.With(zap.String("account_id", id))

	acc := &account.Account{
		Id:          id,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CreatedAt:   timestamppb.Now(),
	}

	err = s.validateNewAccount(ctxlog.ContextWithLogger(ctx, logger), acc)
	if err != nil {
		return nil, err
	}

	err = s.accounts.Create(acc)
	if err != nil {
		logger.Error("Internal Error on AccountDatabase.Create", zap.Error(err))
		return nil, ErrorInternal(err)
	}
	return &account.CreateAccountResponse{
		Account: acc,
	}, nil
}

func (s *AccountService) GetAccount(_ context.Context, req *account.GetAccountRequest) (*account.GetAccountResponse, error) {
	logger := s.logger.Named("GetAccount").With(zap.String("account_id", req.GetAccountId()))

	acc, err := s.accounts.Get(req.GetAccountId())
	if err == nil {
		return &account.GetAccountResponse{
			Account: acc,
		}, nil
	}

	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}

	logger.Error("Internal Error on AccountDatabase.Get", zap.Error(err))
	return nil, ErrorInternal(err)
}

func (s *AccountService) DeleteAccount(_ context.Context, req *account.DeleteAccountRequest) (*account.DeleteAccountResponse, error) {
	logger := s.logger.Named("DeleteAccount").With(zap.String("account_id", req.GetAccountId()))

	err := s.accounts.Delete(req.GetAccountId())
	if err == nil {
		return &account.DeleteAccountResponse{}, nil
	}

	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}

	logger.Error("Internal Error on AccountDatabase.Delete", zap.Error(err))
	return nil, ErrorInternal(err)
}

func updateAccountFields(acc *account.Account, req *account.UpdateAccountRequest) error {
	paths := req.GetUpdateMask().GetPaths()
	if stringutil.SliceContains(paths, AccountNameFieldName) {
		acc.Name = req.GetName()
	}
	if stringutil.SliceContains(paths, AccountDescFieldName) {
		acc.Description = req.GetDescription()
	}
	return nil
}

func (s *AccountService) UpdateAccount(_ context.Context, req *account.UpdateAccountRequest) (*account.UpdateAccountResponse, error) {
	logger := s.logger.Named("UpdateAccount").
		With(zap.String("account_id", req.GetAccountId())).
		With(zap.Strings("update_mask", req.GetUpdateMask().GetPaths()))

	acc, err := s.accounts.Get(req.GetAccountId())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}
	if err != nil {
		logger.Error("Internal Error on AccountDatabase.Get", zap.Error(err))
		return nil, ErrorInternal(err)
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
	if err != nil {
		logger.Error("Internal Error on AccountDatabase.Update", zap.Error(err))
		return nil, ErrorInternal(err)
	}

	return &account.UpdateAccountResponse{
		Account: acc,
	}, nil
}

func (s *AccountService) GetAccountBySshKey(_ context.Context, req *account.GetAccountBySshKeyRequest) (*account.GetAccountBySshKeyResponse, error) {
	logger := s.logger.Named("GetAccountBySshKey")

	sshPubKey, err := sshutil.ParsePublicKey([]byte(req.GetSshPubKey()))
	if err != nil {
		logger.Debug("Ssh Public key parsing failed", zap.String("public key", req.GetSshPubKey()))
		return nil, status.Error(codes.InvalidArgument, "Invalid ssh public key")
	}
	fingerprint := sshutil.GetFingerprint(sshPubKey)
	key, err := s.sshKeys.Get(fingerprint)
	if err != nil {
		if errors.Is(err, database.ErrNotExist{}) {
			return nil, ErrorNotFound("ssh key", fingerprint)
		}

		logger.With(zap.String("fingerprint", fingerprint)).Error("sshKeys.Get Internal Error", zap.Error(err))
		return nil, ErrorInternal(err)
	}

	return &account.GetAccountBySshKeyResponse{
		AccountId: key.GetAccountId(),
	}, nil
}
