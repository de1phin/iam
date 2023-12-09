package service

import (
	"context"

	"github.com/de1phin/iam/services/access/internal/core"
)

//go:generate mockgen -source=storage.go -destination=storage_mock.go -package=service
type Storage interface {
	AddRole(ctx context.Context, role core.Role) error
	GetRole(ctx context.Context, name string) (role core.Role, err error)
	DeleteRole(ctx context.Context, name string) error

	AddAccessBinding(ctx context.Context, binding core.AccessBinding) error
	HaveAccessBinding(ctx context.Context, accountID string, resource string, permission string) (bool, error)
	DeleteAccessBinding(ctx context.Context, binding core.AccessBinding) error
}
