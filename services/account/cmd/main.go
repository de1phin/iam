package main

import (
	"time"

	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/server"
	"github.com/de1phin/iam/services/account/internal/service"
)

func main() {

	accountCache := cache.NewAccountCache()
	sshKeysCache := cache.NewSshKeyCache()

	service := service.NewAccountService(
		service.AccountDatabase(accountCache),
		service.SshKeyDatabase(sshKeysCache),
	)

	cfg := &server.AccountServiceServerConfig{
		AccountService: service,
		// logger
		Address:           ":8443",
		ConnectionTimeout: time.Second * 30,
		// *tls.Certificate
	}

	srv, _ := cfg.RunServer()

	time.Sleep(time.Minute)
	srv.Shutdown()

}
