package client

import (
	"context"

	pb "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type Client interface {
	Authenticate(ctx context.Context, in *pb.AuthenticateRequest, opts ...grpc.CallOption) (*pb.AuthenticateResponse, error)
}

type AccountWrapper struct {
	cl Client
}

func NewAccountWrapper(cl Client) *AccountWrapper {
	return &AccountWrapper{
		cl: cl,
	}
}

func (w *AccountWrapper) Authenticate(ctx context.Context, ssh []byte) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/wrapper/Authenticate")
	defer span.Finish()

	req := &pb.AuthenticateRequest{
		SshKey: ssh,
	}

	resp, err := w.cl.Authenticate(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.AccountId, nil
}
