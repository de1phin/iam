package server

import "github.com/de1phin/iam/services/account/internal/database"

type Server struct {
	accounts database.AccountDatabase
	sshKeys  database.SshKeyDatabase
}

type ServerOpt func(*Server)

func NewServer(opts ...ServerOpt) *Server {
	server := &Server{}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func AccountDatabase(db database.AccountDatabase) ServerOpt {
	return func(s *Server) {
		s.accounts = db
	}
}

func SshKeyDatabase(db database.SshKeyDatabase) ServerOpt {
	return func(s *Server) {
		s.sshKeys = db
	}
}
