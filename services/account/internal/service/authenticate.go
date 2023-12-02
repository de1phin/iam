package service

import (
	"bytes"
	"context"
	"errors"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/database"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AccountService) Authenticate(_ context.Context, req *account.AuthenticateRequest) (*account.AuthenticateResponse, error) {
	signer, err := ssh.ParsePrivateKey(req.GetSshKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Failed to parse ssh key")
	}

	fingerprint := GetFingerprint(signer.PublicKey())

	sshKey, err := s.sshKeys.Get(fingerprint)
	if errors.Is(err, database.ErrNotExist{}) ||
		!bytes.Equal(signer.PublicKey().Marshal(), sshKey.GetPublicKey()) {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &account.AuthenticateResponse{
		AccountId: sshKey.GetAccountId(),
	}, nil
}
