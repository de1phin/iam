package service_test

import (
	"context"
	"testing"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

var (
	sshKey    []byte
	sshPubKey ssh.PublicKey
)

func init() {
	sshKey, sshPubKey = mustGenerateSshKey()
}

func TestInvalidSshKey(t *testing.T) {
	s := service.NewAccountService()

	resp, err := s.Authenticate(context.Background(), &account.AuthenticateRequest{
		SshKey: []byte("bibaboba"),
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "InvalidArgument")
}

func TestUnauthenticated(t *testing.T) {
	s := service.NewAccountService(
		service.AccountDatabase(cache.NewAccountCache()),
		service.SshKeyDatabase(cache.NewSshKeyCache()),
	)

	resp, err := s.Authenticate(context.Background(), &account.AuthenticateRequest{
		SshKey: []byte(sshKey),
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "Unauthenticated")
}

func TestOK(t *testing.T) {
	accounts := cache.NewAccountCache()
	sshKeys := cache.NewSshKeyCache()
	s := service.NewAccountService(
		service.AccountDatabase(accounts),
		service.SshKeyDatabase(sshKeys),
	)

	const accountId = "tima15654655"
	assert.NoError(t, accounts.Create(&account.Account{
		Id: accountId,
	}))
	assert.NoError(t, sshKeys.Create(&account.SshKey{
		AccountId:   accountId,
		Fingerprint: service.GetFingerprint(sshPubKey),
		PublicKey:   sshPubKey.Marshal(),
	}))

	resp, err := s.Authenticate(context.Background(), &account.AuthenticateRequest{
		SshKey: sshKey,
	})
	assert.NoError(t, err)
	assert.Equal(t, resp.GetAccountId(), accountId)
}
