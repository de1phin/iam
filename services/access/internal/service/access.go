package service

import (
	"context"

	"github.com/de1phin/iam/services/access/internal/core"
)

type AccessService struct {
	store          Storage
	tokenValidator TokenValidator
}

func New(store Storage, tokenValidator TokenValidator) *AccessService {
	return &AccessService{
		store:          store,
		tokenValidator: tokenValidator,
	}
}

func (s *AccessService) AddRole(ctx context.Context, role core.Role) error {
	return s.store.AddRole(ctx, role)
}

func (s *AccessService) GetRole(ctx context.Context, name string) (role core.Role, err error) {
	return s.store.GetRole(ctx, name)
}

func (s *AccessService) DeleteRole(ctx context.Context, name string) error {
	return s.store.DeleteRole(ctx, name)
}

func (s *AccessService) AddAccessBinding(ctx context.Context, binding core.AccessBinding) error {
	return s.store.AddAccessBinding(ctx, binding)
}

func (s *AccessService) CheckPermission(ctx context.Context, token string, resource string, permission string) (bool, error) {
	userID, isValid, err := s.tokenValidator.ValidateToken(ctx, token)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, core.ErrInvalidToken
	}

	return s.store.HaveAccessBinding(ctx, userID, resource, permission)
}

func (s *AccessService) DeleteAccessBinding(ctx context.Context, binding core.AccessBinding) error {
	return s.store.DeleteAccessBinding(ctx, binding)
}
