package service

import (
	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/database"
	"go.uber.org/zap"
)

type AccountService struct {
	accounts database.AccountDatabase
	sshKeys  database.SshKeyDatabase

	logger *zap.Logger

	account.UnimplementedAccountServiceServer
}

type AccountServiceOpt func(*AccountService)

func NewAccountService(opts ...AccountServiceOpt) *AccountService {
	server := &AccountService{}

	for _, opt := range opts {
		opt(server)
	}

	if server.logger == nil {
		server.logger = zap.NewNop()
	}

	return server
}

func AccountDatabase(db database.AccountDatabase) AccountServiceOpt {
	return func(s *AccountService) {
		s.accounts = db
	}
}

func SshKeyDatabase(db database.SshKeyDatabase) AccountServiceOpt {
	return func(s *AccountService) {
		s.sshKeys = db
	}
}

func Logger(l *zap.Logger) AccountServiceOpt {
	return func(as *AccountService) {
		as.logger = l
	}
}
