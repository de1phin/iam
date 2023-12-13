package grpcclient

import (
	"context"

	tokenapi "github.com/de1phin/iam/genproto/services/token/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type TokenServiceClient struct {
	api tokenapi.TokenServiceClient
}

func NewTokenServiceClient(
	ctx context.Context,
	addr string,
) (*TokenServiceClient, error) {
	cc, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &TokenServiceClient{
		api: tokenapi.NewTokenServiceClient(cc),
	}, nil
}

func (c *TokenServiceClient) ValidateToken(ctx context.Context, token string) (accountID string, isValid bool, err error) {
	request := &tokenapi.ExchangeTokenRequest{Token: token}
	response, err := c.api.ExchangeToken(ctx, request)

	status := status.Convert(err)
	if status != nil && status.Code() == codes.NotFound {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}
	return response.GetAccountId(), true, nil
}
