package facade

import (
	"context"
	"testing"

	"github.com/de1phin/iam/services/token/internal/facade/mocks"
	"github.com/de1phin/iam/services/token/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	ssh   = "AHAHHAHAHAHHAHAHA"
	token = "APHPPHPHPPHPPPAPFPAGPAGPGPAPHPP"
)

var modelToken = &model.Token{Token: token}

func Test_GenerateToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expect    *model.Token
		expectErr error

		setTokenErr error
	}{
		{
			name:      "ok",
			expect:    convertToModelToken(token),
			expectErr: nil,

			setTokenErr: nil,
		},
		{
			name:      "SetToken Error",
			expect:    nil,
			expectErr: assert.AnError,

			setTokenErr: assert.AnError,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				repo      = &mocks.Repository{}
				generator = &mocks.Generator{}
			)

			generator.EXPECT().Generate().Return(token)
			repo.EXPECT().SetToken(mock.Anything, ssh, mock.Anything).Return(test.setTokenErr)

			facade := NewFacade(&mocks.Cache{}, repo, generator, true)

			actual, err := facade.GenerateToken(context.Background(), ssh)
			assert.Equal(t, test.expect, actual)
			assert.ErrorIs(t, test.expectErr, err)
		})
	}
}

func Test_onlyCache_RefreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		onlyCache bool
		expect    *model.Token
		expectErr error

		setTokenErr error
		setSshErr   error
	}{
		{
			name:      "default",
			onlyCache: true,
			expect:    modelToken,
		},
		{
			name:        "error in SetToken",
			onlyCache:   true,
			expect:      nil,
			setTokenErr: assert.AnError,
			expectErr:   assert.AnError,
		},
		{
			name:      "error in SetSsh",
			onlyCache: true,
			expect:    nil,
			setSshErr: assert.AnError,
			expectErr: assert.AnError,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				cache     = &mocks.Cache{}
				generator = &mocks.Generator{}
			)

			generator.EXPECT().Generate().Return(token)
			cache.EXPECT().SetToken(mock.Anything, ssh, token).Return(test.setTokenErr)
			cache.EXPECT().SetSsh(mock.Anything, token, ssh).Return(test.setSshErr)

			facade := NewFacade(cache, &mocks.Repository{}, generator, test.onlyCache)

			actual, err := facade.RefreshToken(context.Background(), ssh)
			assert.Equal(t, test.expect, actual)
			assert.ErrorIs(t, test.expectErr, err)
		})
	}
}

func Test_onlyCache_DeleteToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		onlyCache bool
		expectErr error

		getToken       string
		getTokenErr    error
		deleteTokenErr error
		deleteSshErr   error
	}{
		{
			name:      "default",
			onlyCache: true,
			getToken:  token,
		},
		{
			name:        "error in GetToken",
			onlyCache:   true,
			getTokenErr: assert.AnError,
			expectErr:   assert.AnError,
		},
		{
			name:           "error in DeleteToken",
			onlyCache:      true,
			getToken:       token,
			deleteTokenErr: assert.AnError,
			expectErr:      assert.AnError,
		},
		{
			name:         "error in DeleteSsh",
			onlyCache:    true,
			getToken:     token,
			deleteSshErr: assert.AnError,
			expectErr:    assert.AnError,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var cache = &mocks.Cache{}

			cache.EXPECT().GetToken(mock.Anything, ssh).Return(test.getToken, test.getTokenErr)
			cache.EXPECT().DeleteToken(mock.Anything, ssh).Return(test.deleteTokenErr)
			cache.EXPECT().DeleteSsh(mock.Anything, token).Return(test.deleteSshErr)

			facade := NewFacade(cache, &mocks.Repository{}, &mocks.Generator{}, test.onlyCache)

			err := facade.DeleteToken(context.Background(), ssh)
			assert.ErrorIs(t, test.expectErr, err)
		})
	}
}

func Test_onlyCache_GetSshByToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		onlyCache bool
		expect    string
		expectErr error

		getSsh    string
		getSshErr error
	}{
		{
			name:      "default",
			onlyCache: true,
			expect:    ssh,
			getSsh:    ssh,
		},
		{
			name:      "error in GetSsh",
			onlyCache: true,
			getSsh:    "",
			expect:    "",
			getSshErr: assert.AnError,
			expectErr: assert.AnError,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var cache = &mocks.Cache{}

			cache.EXPECT().GetSsh(mock.Anything, token).Return(test.getSsh, test.getSshErr)

			facade := NewFacade(cache, &mocks.Repository{}, &mocks.Generator{}, test.onlyCache)

			actual, err := facade.GetSshByToken(context.Background(), *modelToken)
			assert.Equal(t, test.expect, actual)
			assert.ErrorIs(t, test.expectErr, err)
		})
	}
}
