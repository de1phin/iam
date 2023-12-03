package service

import (
	"context"
	"errors"
	"fmt"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/database"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetFingerprint(pubKey ssh.PublicKey) string {
	return ssh.FingerprintSHA256(pubKey)
}

func (s *AccountService) CreateSshKey(_ context.Context, req *account.CreateSshKeyRequest) (*account.CreateSshKeyResponse, error) {
	_, err := s.accounts.Get(req.GetAccountId())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}
	if err != nil {
		return nil, ErrorInternal()
	}

	pubKey, err := ssh.ParsePublicKey(req.GetPublicKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "ssh public key is unparseable")
	}

	fingerprint := GetFingerprint(pubKey)

	// check if ssh key is already used
	_, err = s.sshKeys.Get(fingerprint)
	if err == nil {
		return nil, ErrorAlreadyExists("ssh key")
	}
	if err != nil && !errors.Is(err, database.ErrNotExist{}) {
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
		return nil, ErrorInternal()
	}

	return &account.CreateSshKeyResponse{
		Key: key,
	}, nil
}

func (s *AccountService) ListSshKeys(_ context.Context, req *account.ListSshKeysRequest) (*account.ListSshKeysResponse, error) {
	_, err := s.accounts.Get(req.GetAccountId())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorAccountNotFound(req.GetAccountId())
	}
	keys, err := s.sshKeys.List(req.GetAccountId())
	// not exist is a valid response with 0 ssh keys
	if err != nil && !errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorInternal()
	}

	return &account.ListSshKeysResponse{
		Keys: keys,
	}, nil
}

func (s *AccountService) DeleteSshKey(_ context.Context, req *account.DeleteSshKeyRequest) (*account.DeleteSshKeyResponse, error) {
	err := s.sshKeys.Delete(req.GetAccountId(), req.GetKeyFingerprint())
	if errors.Is(err, database.ErrNotExist{}) {
		return nil, ErrorNotFound("ssh key", fmt.Sprintf("%s/%s", req.GetAccountId(), req.GetKeyFingerprint()))
	}
	if err != nil {
		return nil, ErrorInternal()
	}

	return &account.DeleteSshKeyResponse{}, nil
}
