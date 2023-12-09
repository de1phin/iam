package service

import (
	"context"
	"testing"

	"github.com/de1phin/iam/services/access/internal/core"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mock "github.com/de1phin/iam/services/access/internal/service/mock"
)

type testEnv struct {
	storageMock        *mock.MockStorage
	tokenValidatorMock *mock.MockTokenValidator
	service            *AccessService
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	ctrl := gomock.NewController(t)
	storageMock := mock.NewMockStorage(ctrl)
	tokenValidatorMock := mock.NewMockTokenValidator(ctrl)
	service := New(storageMock, tokenValidatorMock)

	return &testEnv{
		storageMock:        storageMock,
		tokenValidatorMock: tokenValidatorMock,
		service:            service,
	}
}

func TestAccessService_CheckPermission(t *testing.T) {
	env := newTestEnv(t)
	t.Run("ValidTokenWithPermissions", func(t *testing.T) {
		expectedToken := "validtoken"
		expectedPermission := "EatBatat"
		expectedResource := "panfilova"

		expectedAccountID := "khomyak"
		env.tokenValidatorMock.EXPECT().ValidateToken(gomock.Any(), expectedToken).Return(expectedAccountID, true, nil)
		env.storageMock.EXPECT().HaveAccessBinding(
			gomock.Any(),
			expectedAccountID,
			expectedResource,
			expectedPermission,
		).Return(
			true,
			nil,
		)

		have, err := env.service.CheckPermission(context.Background(), expectedToken, expectedResource, expectedPermission)
		require.NoError(t, err)
		require.True(t, have)
	})
	t.Run("ValidTokenWithoutPermissions", func(t *testing.T) {
		expectedToken := "validtoken"
		expectedPermission := "EatBatat"
		expectedResource := "panfilova"

		expectedAccountID := "khomyak"
		env.tokenValidatorMock.EXPECT().ValidateToken(gomock.Any(), expectedToken).Return(expectedAccountID, true, nil)
		env.storageMock.EXPECT().HaveAccessBinding(
			gomock.Any(),
			expectedAccountID,
			expectedResource,
			expectedPermission,
		).Return(
			false,
			nil,
		)

		have, err := env.service.CheckPermission(context.Background(), expectedToken, expectedResource, expectedPermission)
		require.NoError(t, err)
		require.False(t, have)
	})
	t.Run("InvalidToken", func(t *testing.T) {
		expectedToken := "invalidToken"

		env.tokenValidatorMock.EXPECT().ValidateToken(gomock.Any(), expectedToken).Return("", false, nil)

		_, err := env.service.CheckPermission(context.Background(), expectedToken, "", "")
		require.ErrorIs(t, err, core.ErrInvalidToken)
	})

}
