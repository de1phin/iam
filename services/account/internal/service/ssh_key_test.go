package service_test

import (
	"context"
	"testing"
	"time"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/pkg/sshutil"
	"github.com/de1phin/iam/services/account/internal/database"
	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvalidSshKey(t *testing.T) {
	accounts := cache.NewAccountCache()
	accoundService := service.NewAccountService(
		service.AccountDatabase(accounts),
		service.SshKeyDatabase(cache.NewSshKeyCache()),
	)

	const accoundId = "account_id"
	assert.NoError(t, accounts.Create(&account.Account{
		Id: accoundId,
	}))

	resp, err := accoundService.CreateSshKey(context.Background(), &account.CreateSshKeyRequest{
		AccountId: accoundId,
		PublicKey: "invalid public key",
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "InvalidArgument")
	assert.Nil(t, resp)
}

func TestCreateSshKeyAccountNotFound(t *testing.T) {
	accountService := service.NewAccountService(
		service.AccountDatabase(cache.NewAccountCache()),
		service.SshKeyDatabase(cache.NewSshKeyCache()),
	)

	_, sshPubKey := mustGenerateSshKey()
	_, err := accountService.CreateSshKey(context.Background(), &account.CreateSshKeyRequest{
		AccountId: "account_id",
		PublicKey: string(sshPubKey.Marshal()),
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "NotFound")
}

func TestDuplicateSshKey(t *testing.T) {
	accounts := cache.NewAccountCache()
	accountService := service.NewAccountService(
		service.AccountDatabase(accounts),
		service.SshKeyDatabase(cache.NewSshKeyCache()),
	)

	const accoundId = "account_id"
	assert.NoError(t, accounts.Create(&account.Account{
		Id: accoundId,
	}))

	_, sshPubKey := mustGenerateSshKey()
	createResp, err := accountService.CreateSshKey(context.Background(), &account.CreateSshKeyRequest{
		AccountId: accoundId,
		PublicKey: string(sshPubKey.Marshal()),
	})
	assert.NoError(t, err)
	assert.Equal(t, accoundId, createResp.GetKey().GetAccountId())
	assert.Equal(t, string(sshPubKey.Marshal()), createResp.GetKey().GetPublicKey())
	assert.Equal(t, sshutil.GetFingerprint(sshPubKey), createResp.GetKey().GetFingerprint())
	// check created_at is valid [-1s;now]
	assert.True(t, createResp.GetKey().GetCreatedAt().IsValid())
	assert.True(t, createResp.GetKey().GetCreatedAt().AsTime().Before(time.Now()))
	assert.True(t, createResp.GetKey().GetCreatedAt().AsTime().After(time.Now().Add(-time.Second)))

	_, err = accountService.CreateSshKey(context.Background(), &account.CreateSshKeyRequest{
		AccountId: accoundId,
		PublicKey: string(sshPubKey.Marshal()),
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "AlreadyExists")
}

func TestEmptyListSshKeys(t *testing.T) {
	ctx := context.Background()
	accounts := cache.NewAccountCache()
	accountService := service.NewAccountService(
		service.AccountDatabase(accounts),
		service.SshKeyDatabase(cache.NewSshKeyCache()),
	)

	const accountId = "kek"
	// if account does not exist => error
	_, err := accountService.ListSshKeys(ctx, &account.ListSshKeysRequest{
		AccountId: accountId,
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "NotFound")

	assert.NoError(t, accounts.Create(&account.Account{
		Id: accountId,
	}))
	listResp, err := accountService.ListSshKeys(ctx, &account.ListSshKeysRequest{
		AccountId: accountId,
	})
	assert.NoError(t, err)
	assert.NotNil(t, listResp)
	assert.Empty(t, listResp.GetKeys())
}

func TestSshKeyDelete(t *testing.T) {
	accounts := cache.NewAccountCache()
	sshKeys := cache.NewSshKeyCache()
	accountService := service.NewAccountService(
		service.AccountDatabase(accounts),
		service.SshKeyDatabase(sshKeys),
	)

	const accoundId = "boba"
	const fingerprint = "borislav"

	// both account and fingerprint absent
	deleteReq := &account.DeleteSshKeyRequest{
		AccountId:      accoundId,
		KeyFingerprint: fingerprint,
	}
	deleteResp, err := accountService.DeleteSshKey(context.Background(), deleteReq)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "NotFound")
	assert.Nil(t, deleteResp)

	assert.NoError(t, accounts.Create(&account.Account{
		Id: accoundId,
	}))
	assert.NoError(t, sshKeys.Create(&account.SshKey{
		AccountId:   accoundId,
		Fingerprint: fingerprint,
	}))

	deleteResp, err = accountService.DeleteSshKey(context.Background(), deleteReq)
	assert.NoError(t, err)
	assert.NotNil(t, deleteResp)

	_, err = sshKeys.Get(fingerprint)
	assert.ErrorIs(t, err, database.ErrNotExist{})

	deleteResp, err = accountService.DeleteSshKey(context.Background(), deleteReq)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "NotFound")
	assert.Nil(t, deleteResp)
}

func TestSameSshKeyProhibited(t *testing.T) {
	const (
		accountId1 = "1"
		accountId2 = "2"
	)

	_, sshPubKey := mustGenerateSshKey()

	ctx := context.Background()

	accounts := cache.NewAccountCache()
	accountService := service.NewAccountService(
		service.AccountDatabase(accounts),
		service.SshKeyDatabase(cache.NewSshKeyCache()),
	)

	assert.NoError(t, accounts.Create(&account.Account{
		Id: accountId1,
	}))
	assert.NoError(t, accounts.Create(&account.Account{
		Id: accountId2,
	}))

	_, err := accountService.CreateSshKey(ctx, &account.CreateSshKeyRequest{
		AccountId: accountId1,
		PublicKey: string(sshPubKey.Marshal()),
	})
	assert.NoError(t, err)

	// same account same ssh key prohibited
	_, err = accountService.CreateSshKey(ctx, &account.CreateSshKeyRequest{
		AccountId: accountId1,
		PublicKey: string(sshPubKey.Marshal()),
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "AlreadyExists")

	// different account same ssh key also prohibited
	_, err = accountService.CreateSshKey(ctx, &account.CreateSshKeyRequest{
		AccountId: accountId2,
		PublicKey: string(sshPubKey.Marshal()),
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "AlreadyExists")
}

func TestSshKeyValidCRUD(t *testing.T) {
	ctx := context.Background()
	accounts := cache.NewAccountCache()
	accountService := service.NewAccountService(
		service.AccountDatabase(accounts),
		service.SshKeyDatabase(cache.NewSshKeyCache()),
	)

	const accoundId = "nazim_malyshev"
	assert.NoError(t, accounts.Create(&account.Account{
		Id: accoundId,
	}))

	_, sshPubKey1 := mustGenerateSshKey()

	createResp, err := accountService.CreateSshKey(ctx, &account.CreateSshKeyRequest{
		AccountId: accoundId,
		PublicKey: string(sshPubKey1.Marshal()),
	})
	assert.NoError(t, err)
	assert.Equal(t, accoundId, createResp.GetKey().GetAccountId())
	assert.Equal(t, string(sshPubKey1.Marshal()), createResp.GetKey().GetPublicKey())
	assert.Equal(t, sshutil.GetFingerprint(sshPubKey1), createResp.GetKey().GetFingerprint())
	// check created_at is valid [-1s;now]
	assert.True(t, createResp.GetKey().GetCreatedAt().IsValid())
	assert.True(t, createResp.GetKey().GetCreatedAt().AsTime().Before(time.Now()))
	assert.True(t, createResp.GetKey().GetCreatedAt().AsTime().After(time.Now().Add(-time.Second)))

	_, sshPubKey2 := mustGenerateSshKey()

	createResp, err = accountService.CreateSshKey(ctx, &account.CreateSshKeyRequest{
		AccountId: accoundId,
		PublicKey: string(sshPubKey2.Marshal()),
	})
	assert.NoError(t, err)
	assert.Equal(t, accoundId, createResp.GetKey().GetAccountId())
	assert.Equal(t, string(sshPubKey2.Marshal()), createResp.GetKey().GetPublicKey())
	assert.Equal(t, sshutil.GetFingerprint(sshPubKey2), createResp.GetKey().GetFingerprint())
	// check created_at is valid [-1s;now]
	assert.True(t, createResp.GetKey().GetCreatedAt().IsValid())
	assert.True(t, createResp.GetKey().GetCreatedAt().AsTime().Before(time.Now()))
	assert.True(t, createResp.GetKey().GetCreatedAt().AsTime().After(time.Now().Add(-time.Second)))

	listResp, err := accountService.ListSshKeys(ctx, &account.ListSshKeysRequest{
		AccountId: accoundId,
	})
	assert.NoError(t, err)
	assert.Len(t, listResp.GetKeys(), 2)

	assert.Equal(t, listResp.GetKeys()[0].GetAccountId(), accoundId)
	assert.Equal(t, listResp.GetKeys()[0].GetFingerprint(), sshutil.GetFingerprint(sshPubKey1))
	assert.Equal(t, listResp.GetKeys()[0].GetPublicKey(), string(sshPubKey1.Marshal()))

	assert.Equal(t, listResp.GetKeys()[1].GetAccountId(), accoundId)
	assert.Equal(t, listResp.GetKeys()[1].GetFingerprint(), sshutil.GetFingerprint(sshPubKey2))
	assert.Equal(t, listResp.GetKeys()[1].GetPublicKey(), string(sshPubKey2.Marshal()))

	_, err = accountService.DeleteSshKey(ctx, &account.DeleteSshKeyRequest{
		AccountId:      accoundId,
		KeyFingerprint: sshutil.GetFingerprint(sshPubKey1),
	})
	assert.NoError(t, err)

	listResp, err = accountService.ListSshKeys(ctx, &account.ListSshKeysRequest{
		AccountId: accoundId,
	})
	assert.NoError(t, err)
	assert.Len(t, listResp.GetKeys(), 1)

	assert.Equal(t, listResp.GetKeys()[0].GetAccountId(), accoundId)
	assert.Equal(t, listResp.GetKeys()[0].GetFingerprint(), sshutil.GetFingerprint(sshPubKey2))
	assert.Equal(t, listResp.GetKeys()[0].GetPublicKey(), string(sshPubKey2.Marshal()))

	_, err = accountService.DeleteSshKey(ctx, &account.DeleteSshKeyRequest{
		AccountId:      accoundId,
		KeyFingerprint: sshutil.GetFingerprint(sshPubKey2),
	})
	assert.NoError(t, err)

	listResp, err = accountService.ListSshKeys(ctx, &account.ListSshKeysRequest{
		AccountId: accoundId,
	})
	assert.NoError(t, err)
	assert.Len(t, listResp.GetKeys(), 0)
}
