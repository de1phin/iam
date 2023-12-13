package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/de1phin/iam/services/access/internal/core"
	"github.com/lib/pq"
)

func (s *Storage) AddRole(ctx context.Context, role core.Role) error {
	addRoleSQL := `INSERT INTO roles(name, permissions) VALUES ($1, $2)`
	_, err := s.conn.ExecContext(ctx, addRoleSQL, role.Name, pq.Array(role.Permissions))
	return err
}

func (s *Storage) GetRole(ctx context.Context, name string) (role core.Role, err error) {
	role.Name = name
	getRoleSQL := `SELECT permissions FROM roles WHERE name = $1`
	err = s.conn.QueryRowContext(ctx, getRoleSQL, name).Scan(pq.Array(&role.Permissions))
	if errors.Is(err, sql.ErrNoRows) {
		return core.Role{}, ErrNotExist{}
	}
	if err != nil {
		return core.Role{}, err
	}

	return role, nil
}

func (s *Storage) DeleteRole(ctx context.Context, name string) error {
	deleteRoleSQL := `DELETE FROM roles WHERE name = $1`
	_, err := s.conn.ExecContext(ctx, deleteRoleSQL, name)
	return err
}
