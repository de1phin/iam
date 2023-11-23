package token

import (
	"context"

	desc "github.com/de1phin/iam/genproto/services/token/api"

	"github.com/de1phin/iam/token/internal/repository"
)

type TokenRepository interface {
}

type Implementation struct {
	tokenRepo TokenRepository

	desc.UnimplementedTokenServiceServer
}

func NewService(tokenRepo repository.Database) *Implementation {
	return &Implementation{
		tokenRepo: tokenRepo,
	}
}

func (i *Implementation) GenerateToken(ctx context.Context, req *desc.GenerateTokenRequest) (*desc.GenerateTokenResponse, error) {
	return nil, nil
}

func (i *Implementation) RefreshToken(ctx context.Context, req *desc.RefreshTokenRequest) (*desc.RefreshTokenResponse, error) {
	return nil, nil
}

func (i *Implementation) DeleteToken(ctx context.Context, req *desc.RemoveTokenRequest) (*desc.RemoveTokenResponse, error) {
	return nil, nil
}

func (i *Implementation) CheckToken(context.Context, *desc.CheckTokenRequest) (*desc.CheckTokenResponse, error) {
	return nil, nil
}
