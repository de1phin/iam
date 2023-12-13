package token

import (
	"context"
	"errors"
	"math/rand"
	"time"

	desc "github.com/de1phin/iam/genproto/services/token/api"
	"github.com/de1phin/iam/pkg/logger"
	"github.com/de1phin/iam/pkg/sshutil"
	"github.com/de1phin/iam/services/token/internal/model"
	"github.com/de1phin/iam/services/token/internal/repository"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TokenFacade interface {
	CreateToken(ctx context.Context, accountId string) (*model.Token, error)
	RefreshToken(ctx context.Context, token string) (*model.Token, error)
	DeleteToken(ctx context.Context, token string) error
	ExpireTokens(ctx context.Context) error
	GetToken(ctx context.Context, token string) (*model.Token, error)
}

type AccountService interface {
	GetAccountBySshKey(ctx context.Context, sshPubKey string) (string, error)
}

type Implementation struct {
	token   TokenFacade
	account AccountService

	quitExpire chan struct{}

	desc.UnimplementedTokenServiceServer
}

func NewService(token TokenFacade, account AccountService) *Implementation {
	return &Implementation{
		token:   token,
		account: account,
	}
}

func (i *Implementation) RunExpirationRoutine(ctx context.Context, intervalMin time.Duration, intervalMax time.Duration) {
	i.quitExpire = make(chan struct{})

	for {
		diff := int64(intervalMax) - int64(intervalMin)
		interval := time.Duration(int64(intervalMin) + rand.Int63n(diff))

		timer := time.NewTimer(interval)

		select {
		case <-i.quitExpire:
			timer.Stop()
			break
		case <-timer.C:
			err := i.token.ExpireTokens(ctx)
			if err != nil {
				logger.Error("ExpireTokens", zap.Error(err))
			}
			break
		}
	}
}

func (i *Implementation) StopExpirationRoutine() {
	if i.quitExpire != nil {
		i.quitExpire <- struct{}{}
		i.quitExpire = nil
	}
}

func (i *Implementation) CreateToken(ctx context.Context, req *desc.CreateTokenRequest) (*desc.CreateTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/GenerateToken")
	defer span.Finish()

	accountId, err := i.account.GetAccountBySshKey(ctx, req.GetSshPubKey())
	reqStatus := status.Convert(err)
	if reqStatus != nil && reqStatus.Code() == codes.InvalidArgument {
		return nil, err
	}
	if reqStatus != nil && reqStatus.Code() == codes.NotFound {
		return nil, status.Error(codes.Unauthenticated, "Unable to authenticate account")
	}

	token, err := i.token.CreateToken(ctx, accountId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	encryptedToken, err := sshutil.EncryptWithPublicKey([]byte(token.Token), []byte(req.GetSshPubKey()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.CreateTokenResponse{
		Token:     string(encryptedToken),
		ExpiresAt: timestamppb.New(token.ExpiresAt),
	}, nil
}

func (i *Implementation) RefreshToken(ctx context.Context, req *desc.RefreshTokenRequest) (*desc.RefreshTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/RefreshToken")
	defer span.Finish()

	token, err := i.token.RefreshToken(ctx, req.GetToken())
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Token not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.RefreshTokenResponse{
		ExpiresAt: timestamppb.New(token.ExpiresAt),
	}, nil
}

func (i *Implementation) DeleteToken(ctx context.Context, req *desc.DeleteTokenRequest) (*desc.DeleteTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/DeleteToken")
	defer span.Finish()

	err := i.token.DeleteToken(ctx, req.GetToken())
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Token not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.DeleteTokenResponse{}, nil
}

func (i *Implementation) ExchangeToken(ctx context.Context, req *desc.ExchangeTokenRequest) (*desc.ExchangeTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api/ExchangeToken")
	defer span.Finish()

	token, err := i.token.GetToken(ctx, req.GetToken())
	if errors.Is(err, repository.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "Token not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.ExchangeTokenResponse{
		AccountId: token.AccountId,
	}, nil
}
