package postgres

import (
	"context"

	"github.com/de1phin/iam/services/access/internal/core"
)

func (s *Storage) AddAccessBinding(ctx context.Context, binding core.AccessBinding) error {
	addAccessBindingSQL := `INSERT INTO access_bindings(account_id, role_name, resource) VALUES ($1, $2, $3)`
	_, err := s.conn.QueryContext(ctx, addAccessBindingSQL, binding.AccountID, binding.RoleName, binding.Resource)
	return err
}

func (s *Storage) HaveAccessBinding(ctx context.Context, accountID string, resource string, permission string) (bool, error) {
	haveAccessBindingSQL := `
	SELECT COUNT(*)
	FROM roles r
	JOIN access_bindings ab ON ab.role_name = r.name
	WHERE ab.account_id = $1
		AND ab.resource = $2
		AND $3 = ANY (r.permissions)
	LIMIT 1
	`

	var count int
	err := s.conn.QueryRowContext(ctx, haveAccessBindingSQL, accountID, resource, permission).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, err
}

func (s *Storage) DeleteAccessBinding(ctx context.Context, binding core.AccessBinding) error {
	deleteAccessBindingSQL := `
	DELETE FROM access_bindings
	WHERE
	account_id = $1
	AND resource = $2
	AND role_name = $3
	`
	_, err := s.conn.QueryContext(ctx, deleteAccessBindingSQL, binding.AccountID, binding.Resource, binding.RoleName)
	return err
}
