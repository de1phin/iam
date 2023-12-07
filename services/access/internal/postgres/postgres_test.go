//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/de1phin/iam/services/access/internal/core"
	"github.com/stretchr/testify/require"
)

type testEnv struct {
	store  *Storage
	ctx    context.Context
	cancel context.CancelFunc
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	env := &testEnv{}
	env.ctx, env.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	env.store, err = New(env.ctx, Options{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "bibaboba",
		DBName:   "postgres",
	})
	require.NoError(t, err)

	env.cleanup(t)

	return env
}

func (e *testEnv) cleanup(t *testing.T) {
	t.Helper()

	cleanupSQL := `
	DROP TABLE IF EXISTS roles;
	DROP TABLE IF EXISTS access_bindings;
	`
	_, err := e.store.conn.ExecContext(e.ctx, cleanupSQL)
	require.NoError(t, err)
	err = e.store.initSchema(e.ctx)
	require.NoError(t, err)
}

func TestStorage(t *testing.T) {
	t.Run("SimpleRoleCRUD", func(t *testing.T) {
		env := newTestEnv(t)
		expectedRole := core.Role{
			Name:        "biba.admin",
			Permissions: []string{"CreateBiba", "ReadBiba", "UpdateBiba", "DeleteBiba"},
		}
		err := env.store.AddRole(env.ctx, expectedRole)
		require.NoError(t, err)

		actualRole, err := env.store.GetRole(env.ctx, expectedRole.Name)
		require.NoError(t, err)
		require.Equal(t, expectedRole, actualRole)

		err = env.store.DeleteRole(env.ctx, expectedRole.Name)
		require.NoError(t, err)

		_, err = env.store.GetRole(env.ctx, expectedRole.Name)
		require.Error(t, err)
	})
	t.Run("SimpleAccessBindingCRUD", func(t *testing.T) {
		env := newTestEnv(t)
		const (
			expectedPermission   = "GetInBodyStats"
			unexpectedPermission = "EatBatat"
		)
		expectedRole := core.Role{
			Name:        "gym.admin",
			Permissions: []string{expectedPermission, "DoExcersize"},
		}
		err := env.store.AddRole(env.ctx, expectedRole)
		require.NoError(t, err)

		expectedAB := core.AccessBinding{
			UserID:   "khomyak",
			RoleName: expectedRole.Name,
			Resource: "ddx_fitness_shukinskaya",
		}
		err = env.store.AddAccessBinding(env.ctx, expectedAB)
		require.NoError(t, err)
		have, err := env.store.HaveAccessBinding(env.ctx, expectedAB.UserID, expectedAB.Resource, expectedPermission)
		require.NoError(t, err)
		require.True(t, have)

		have, err = env.store.HaveAccessBinding(env.ctx, expectedAB.UserID, expectedAB.Resource, unexpectedPermission)
		require.NoError(t, err)
		require.False(t, have)

	})

}
