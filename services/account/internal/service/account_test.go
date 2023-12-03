package service_test

import (
	"context"
	"testing"
	"time"

	account "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/de1phin/iam/services/account/internal/database/cache"
	"github.com/de1phin/iam/services/account/internal/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestCreateAccountWithInvalidID(t *testing.T) {
	accounts := cache.NewAccountCache()
	accountService := service.NewAccountService(
		service.AccountDatabase(accounts),
	)

	const accountName = "test-account"
	ctx := context.Background()
	// invalid symbols
	resp, err := accountService.CreateAccount(ctx, &account.CreateAccountRequest{
		WellKnownId: "41234@#$%^*",
		Name:        accountName,
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "InvalidArgument")
	assert.Nil(t, resp)

	const accountId = "test.account"
	// valid create
	resp, err = accountService.CreateAccount(ctx, &account.CreateAccountRequest{
		WellKnownId: accountId,
		Name:        accountName,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, accountId, resp.GetAccount().GetId())

	// id already in use
	resp, err = accountService.CreateAccount(ctx, &account.CreateAccountRequest{
		WellKnownId: accountId,
		Name:        accountName,
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "AlreadyExists")
	assert.Nil(t, resp)
}

func TestCreateAccountWithInvalidName(t *testing.T) {
	accountService := service.NewAccountService(
		service.AccountDatabase(cache.NewAccountCache()),
	)

	ctx := context.Background()
	accountName := ""

	resp, err := accountService.CreateAccount(ctx, &account.CreateAccountRequest{
		Name: accountName,
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "InvalidArgument")
	assert.Nil(t, resp)

	nameLen := (service.NameMinLength + service.NameMaxLength) / 2
	for len(accountName) < nameLen {
		accountName += "a"
	}

	resp, err = accountService.CreateAccount(ctx, &account.CreateAccountRequest{
		Name: accountName,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, accountName, resp.GetAccount().GetName())

	nameLen = service.NameMaxLength + 1
	for len(accountName) < nameLen {
		accountName += "a"
	}
	resp, err = accountService.CreateAccount(ctx, &account.CreateAccountRequest{
		Name: accountName,
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "InvalidArgument")
	assert.Nil(t, resp)
}

func TestAccountCRUD(t *testing.T) {
	accounts := cache.NewAccountCache()
	accountService := service.NewAccountService(
		service.AccountDatabase(accounts),
	)

	const (
		accountId   = "ivan_h"
		accountName = "ivakhomiakov"
		accountDesc = "lalalalala"
	)
	ctx := context.Background()

	createAccountRequest := &account.CreateAccountRequest{
		WellKnownId: "",
		Name:        "",
		Description: accountDesc,
	}
	// empty id is not allowed
	_, err := accountService.CreateAccount(ctx, createAccountRequest)
	assert.Error(t, err)
	createAccountRequest.WellKnownId = accountId
	// empty name is not allowed
	_, err = accountService.CreateAccount(ctx, createAccountRequest)
	assert.Error(t, err)
	createAccountRequest.Name = accountName

	createResp, err := accountService.CreateAccount(ctx, createAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, accountId, createResp.GetAccount().GetId())
	assert.Equal(t, accountName, createResp.GetAccount().GetName())
	assert.Equal(t, accountDesc, createResp.GetAccount().GetDescription())

	getAccountRequest := &account.GetAccountRequest{
		AccountId: accountId,
	}
	getResp, err := accountService.GetAccount(ctx, getAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, accountId, getResp.GetAccount().GetId())
	assert.Equal(t, accountName, getResp.GetAccount().GetName())
	assert.Equal(t, accountDesc, getResp.GetAccount().GetDescription())
	// assert created_at is valid ([-1s; now])
	assert.True(t, getResp.GetAccount().GetCreatedAt().IsValid())
	assert.True(t, getResp.GetAccount().GetCreatedAt().AsTime().Before(time.Now()))
	assert.True(t, getResp.GetAccount().GetCreatedAt().AsTime().After(time.Now().Add(-time.Second)))

	_, err = accountService.GetAccount(ctx, &account.GetAccountRequest{AccountId: "ivan.h"})
	assert.Error(t, err)

	const (
		updatedAccountName = "not_ivakhomiakov"
		updatedAccountDesc = ""
	)
	// update both name and desc fields
	updateMask, err := fieldmaskpb.New(&account.UpdateAccountRequest{},
		service.AccountNameFieldName, service.AccountDescFieldName)
	assert.NoError(t, err)
	updateAccountRequest := &account.UpdateAccountRequest{
		AccountId:   "",
		Name:        "",
		Description: updatedAccountDesc,

		UpdateMask: updateMask,
	}
	// account id "" does not exist
	_, err = accountService.UpdateAccount(ctx, updateAccountRequest)
	assert.Error(t, err)
	updateAccountRequest.AccountId = accountId
	// name must not be ""
	_, err = accountService.UpdateAccount(ctx, updateAccountRequest)
	assert.Error(t, err)

	updateAccountRequest.Name = updatedAccountName
	updResp, err := accountService.UpdateAccount(ctx, updateAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, updResp.GetAccount().GetId(), accountId)
	assert.Equal(t, updResp.GetAccount().GetName(), updatedAccountName)
	assert.Equal(t, updResp.GetAccount().GetDescription(), updatedAccountDesc)

	getResp, err = accountService.GetAccount(ctx, getAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, getResp.GetAccount().GetId(), accountId)
	assert.Equal(t, getResp.GetAccount().GetName(), updatedAccountName)
	assert.Equal(t, getResp.GetAccount().GetDescription(), updatedAccountDesc)

	// update back account name
	updateAccountRequest.UpdateMask, err = fieldmaskpb.New(&account.UpdateAccountRequest{},
		service.AccountNameFieldName)
	updateAccountRequest.Name = accountName
	assert.NoError(t, err)
	updResp, err = accountService.UpdateAccount(ctx, updateAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, updResp.GetAccount().GetId(), accountId)
	assert.Equal(t, updResp.GetAccount().GetName(), accountName)
	assert.Equal(t, updResp.GetAccount().GetDescription(), updatedAccountDesc)

	getResp, err = accountService.GetAccount(ctx, getAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, getResp.GetAccount().GetId(), accountId)
	assert.Equal(t, getResp.GetAccount().GetName(), accountName)
	assert.Equal(t, getResp.GetAccount().GetDescription(), updatedAccountDesc)

	// update back account description
	updateAccountRequest.UpdateMask, err = fieldmaskpb.New(&account.UpdateAccountRequest{},
		service.AccountDescFieldName)
	updateAccountRequest.Description = accountDesc
	assert.NoError(t, err)
	updResp, err = accountService.UpdateAccount(ctx, updateAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, updResp.GetAccount().GetId(), accountId)
	assert.Equal(t, updResp.GetAccount().GetName(), accountName)
	assert.Equal(t, updResp.GetAccount().GetDescription(), accountDesc)

	getResp, err = accountService.GetAccount(ctx, getAccountRequest)
	assert.NoError(t, err)
	assert.Equal(t, getResp.GetAccount().GetId(), accountId)
	assert.Equal(t, getResp.GetAccount().GetName(), accountName)
	assert.Equal(t, getResp.GetAccount().GetDescription(), accountDesc)

	_, err = accountService.DeleteAccount(ctx, &account.DeleteAccountRequest{AccountId: accountId})
	assert.NoError(t, err)

	_, err = accountService.GetAccount(ctx, getAccountRequest)
	assert.Error(t, err)

	_, err = accountService.DeleteAccount(ctx, &account.DeleteAccountRequest{AccountId: accountId})
	assert.Error(t, err)
}
