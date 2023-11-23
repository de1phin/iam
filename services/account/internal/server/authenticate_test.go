package server_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"testing"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

var (
	sshKey    []byte
	sshPubKey ssh.PublicKey
)

func init() {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	sshPubKey, err = ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		panic(err)
	}
	block, err := ssh.MarshalPrivateKey(rsaKey, "")
	if err != nil {
		panic(err)
	}
	sshKey = pem.EncodeToMemory(block)
}

func TestInvalidSshKey(t *testing.T) {
	s := server.NewServer()

	resp, err := s.Authenticate(context.Background(), &account.AuthenticateRequest{
		SshKey: []byte("bibaboba"),
	})

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "InvalidArgument")
}

func TestUnauthenticated(t *testing.T) {
	s := server.NewServer(
		server.AccountDatabase(cache.NewAccountCache()),
		server.SshKeyDatabase(cache.NewSshKeyCache()),
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
	s := server.NewServer(
		server.AccountDatabase(accounts),
		server.SshKeyDatabase(sshKeys),
	)

	const accountId = "tima15654655"
	assert.NoError(t, accounts.Create(&account.Account{
		Id: accountId,
	}))
	assert.NoError(t, sshKeys.Create(&account.SshKey{
		AccountId:   accountId,
		Fingerprint: server.GetFingerprint(sshPubKey),
		PublicKey:   sshPubKey.Marshal(),
	}))

	resp, err := s.Authenticate(context.Background(), &account.AuthenticateRequest{
		SshKey: sshKey,
	})
	assert.NoError(t, err)
	assert.Equal(t, resp.GetAccountId(), accountId)
}
