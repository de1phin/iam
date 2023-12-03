package service

import (
	"context"
	"errors"
	"fmt"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/database"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetFingerprint(pubKey ssh.PublicKey) string {
	return ssh.FingerprintSHA256(pubKey)
}

func (s *AccountService) CreateSshKey(_ context.Context, req *account.CreateSshKeyRequest) (*account.CreateSshKeyResponse, error) {
	logger := s.logger.Named("CreateSshKey").With(zap.String("account_id", req.GetAccountId()))

	_, err := s.accounts.Get(req.GetAccountId())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}
	if err != nil {
		logger.Error("Internal Error on AccountDatabase.Get", zap.Error(err))
		return nil, ErrorInternal()
	}

	pubKey, err := ssh.ParsePublicKey(req.GetPublicKey())
	if err != nil {
		logger.Debug("Ssh Public Key parsing failed", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Failed to parse ssh public key")
	}

	fingerprint := GetFingerprint(pubKey)
	logger = logger.With(zap.String("fingerprint", fingerprint))

	// check if ssh key is already used
	_, err = s.sshKeys.Get(fingerprint)
	if err == nil {
		return nil, ErrorAlreadyExists("ssh key")
	}
	if err != nil && !errors.Is(err, database.ErrNotExist{}) {
		logger.Error("Internal Error on SshKeyDatabase.Get", zap.Error(err))
		return nil, ErrorInternal()
	}

	key := &account.SshKey{
		Fingerprint: fingerprint,
		AccountId:   req.GetAccountId(),
		PublicKey:   pubKey.Marshal(),
		CreatedAt:   timestamppb.Now(),
	}
	err = s.sshKeys.Create(key)
	if errors.Is(err, database.ErrAlreadyExists{}) {
		return nil, ErrorAlreadyExists("ssh key")
	}
	if err != nil {
		logger.Error("Internal Error on SshKeyDatabase.Create", zap.Error(err))
		return nil, ErrorInternal()
	}

	return &account.CreateSshKeyResponse{
		Key: key,
	}, nil
}

func (s *AccountService) ListSshKeys(_ context.Context, req *account.ListSshKeysRequest) (*account.ListSshKeysResponse, error) {
	logger := s.logger.Named("ListSshKeys").With(zap.String("account_id", req.GetAccountId()))

	_, err := s.accounts.Get(req.GetAccountId())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}
	keys, err := s.sshKeys.List(req.GetAccountId())
	// not exist is a valid response with 0 ssh keys
	if err != nil && !errors.Is(err, database.ErrNotExist{}) {
		logger.Error("Internal Error on SshKeyDatabase.List", zap.Error(err))
		return nil, ErrorInternal()
	}

	return &account.ListSshKeysResponse{
		Keys: keys,
	}, nil
}

func (s *AccountService) DeleteSshKey(_ context.Context, req *account.DeleteSshKeyRequest) (*account.DeleteSshKeyResponse, error) {
	logger := s.logger.Named("DeleteSshKey").
		With(zap.String("account_id", req.GetAccountId())).
		With(zap.String("fingerprint", req.GetKeyFingerprint()))

	err := s.sshKeys.Delete(req.GetAccountId(), req.GetKeyFingerprint())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorNotFound("ssh key", fmt.Sprintf("%s/%s", req.GetAccountId(), req.GetKeyFingerprint()))
	}
	if err != nil {
		logger.Error("Internal Error on SshKeyDatabase.Delete", zap.Error(err))
		return nil, ErrorInternal()
	}

	return &account.DeleteSshKeyResponse{}, nil
}
