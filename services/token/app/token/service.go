package token

import (
	"context"

	desc "github.com/de1phin/iam/genproto/services/token/api"
	"github.com/de1phin/iam/pkg/sshutil"
	"github.com/de1phin/iam/services/token/internal/mapping"
	"github.com/de1phin/iam/services/token/internal/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenFacade interface {
	GenerateToken(ctx context.Context, ssh string) (*model.Token, error)
	RefreshToken(ctx context.Context, ssh string) error
	DeleteToken(ctx context.Context, ssh string) error
	GetSshByToken(ctx context.Context, tk model.Token) (string, error)
}

type AccountService interface {
	GetAccountBySshKey(ctx context.Context, sshPubKey []byte) (string, error)
}

type Implementation struct {
	token   TokenFacade
	account AccountService

	desc.UnimplementedTokenServiceServer
}

func NewService(token TokenFacade, account AccountService) *Implementation {
	return &Implementation{
		token:   token,
		account: account,
	}
}

func (i *Implementation) GenerateToken(ctx context.Context, req *desc.GenerateTokenRequest) (*desc.GenerateTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/GenerateToken")
	defer span.Finish()

	if len(req.GetSshPubKey()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty ssh")
	}

	token, err := i.token.GenerateToken(ctx, string(req.GetSshPubKey()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	encryptedToken, err := sshutil.EncryptWithPublicKey([]byte(token.Token), []byte(req.GetSshPubKey()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.GenerateTokenResponse{
		Token: &desc.Token{
			Token:  string(encryptedToken),
			Status: desc.TokenStatus_VALID,
		},
	}, nil
}

func (i *Implementation) RefreshToken(ctx context.Context, req *desc.RefreshTokenRequest) (*desc.RefreshTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/RefreshToken")
	defer span.Finish()

	if len(req.GetToken()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty token")
	}

	if err := i.token.RefreshToken(ctx, string(req.GetToken())); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &desc.RefreshTokenResponse{
		Token: &desc.Token{
			Token:  string(req.GetToken()),
			Status: desc.TokenStatus_VALID,
		},
	}, nil
}

func (i *Implementation) DeleteToken(ctx context.Context, req *desc.RemoveTokenRequest) (*desc.RemoveTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/DeleteToken")
	defer span.Finish()

	if len(req.GetSshKey()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty ssh")
	}

	if err := i.token.DeleteToken(ctx, string(req.GetSshKey())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.RemoveTokenResponse{}, nil
}

func (i *Implementation) CheckToken(ctx context.Context, req *desc.CheckTokenRequest) (*desc.CheckTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/CheckToken")
	defer span.Finish()

	if req.GetToken() == nil || req.GetToken().Token == "" {
		return nil, status.Error(codes.InvalidArgument, "empty token")
	}

	token := mapping.MapTokenToInternal(req.GetToken())

	ssh, err := i.token.GetSshByToken(ctx, token)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if ssh == "" {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	userID, err := i.account.GetAccountBySshKey(ctx, []byte(ssh))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.CheckTokenResponse{
		UserId: userID,
	}, nil
}
