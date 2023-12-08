package grpcclient

import (
	"context"

	tokenapi "github.com/de1phin/iam/genproto/services/token/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func (c *TokenServiceClient) ValidateToken(ctx context.Context, token string) (userID string, isValid bool, err error) {
	request := &tokenapi.CheckTokenRequest{Token: &tokenapi.Token{Token: token}}
	response, err := c.api.CheckToken(ctx, request)
	if err != nil {
		return "", false, err
	}
	if response.Status == tokenapi.TokenStatus_UNDEFINED {
		return "", false, nil
	}
	return response.UserId, true, nil
}
