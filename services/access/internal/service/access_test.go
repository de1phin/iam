package service

import (
	"context"
	"testing"

	"github.com/de1phin/iam/services/access/internal/core"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	storageMock        *MockStorage
	tokenValidatorMock *MockTokenValidator
	service            *AccessService
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	ctrl := gomock.NewController(t)
	storageMock := NewMockStorage(ctrl)
	tokenValidatorMock := NewMockTokenValidator(ctrl)
	service := New(storageMock, tokenValidatorMock)

	return &testEnv{
		storageMock:        storageMock,
		tokenValidatorMock: tokenValidatorMock,
		service:            service,
	}
}

func TestAccessService_HaveAccessBinding(t *testing.T) {
	env := newTestEnv(t)
	t.Run("ValidTokenWithPermissions", func(t *testing.T) {
		expectedToken := "validtoken"
		expectedPermission := "EatBatat"
		expectedResource := "panfilova"

		expectedUserID := "khomyak"
		env.tokenValidatorMock.EXPECT().ValidateToken(gomock.Any(), expectedToken).Return(expectedUserID, true, nil)
		env.storageMock.EXPECT().HaveAccessBinding(
			gomock.Any(),
			expectedUserID,
			expectedResource,
			expectedPermission,
		).Return(
			true,
			nil,
		)

		have, err := env.service.HaveAccessBinding(context.Background(), expectedToken, expectedResource, expectedPermission)
		require.NoError(t, err)
		require.True(t, have)
	})
	t.Run("ValidTokenWithoutPermissions", func(t *testing.T) {
		expectedToken := "validtoken"
		expectedPermission := "EatBatat"
		expectedResource := "panfilova"

		expectedUserID := "khomyak"
		env.tokenValidatorMock.EXPECT().ValidateToken(gomock.Any(), expectedToken).Return(expectedUserID, true, nil)
		env.storageMock.EXPECT().HaveAccessBinding(
			gomock.Any(),
			expectedUserID,
			expectedResource,
			expectedPermission,
		).Return(
			false,
			nil,
		)

		have, err := env.service.HaveAccessBinding(context.Background(), expectedToken, expectedResource, expectedPermission)
		require.NoError(t, err)
		require.False(t, have)
	})
	t.Run("InvalidToken", func(t *testing.T) {
		expectedToken := "invalidToken"

		env.tokenValidatorMock.EXPECT().ValidateToken(gomock.Any(), expectedToken).Return("", false, nil)

		_, err := env.service.HaveAccessBinding(context.Background(), expectedToken, "", "")
		require.ErrorIs(t, err, core.ErrInvalidToken)
	})

}
