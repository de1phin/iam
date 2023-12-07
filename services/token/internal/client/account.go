package client

import (
	"context"

	pb "github.com/de1phin/iam/genproto/services/account/api"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type Client interface {
	GetAccountBySshKey(ctx context.Context, in *pb.GetAccountBySshKeyRequest, opts ...grpc.CallOption) (*pb.GetAccountBySshKeyResponse, error)
}

type AccountWrapper struct {
	cl Client
}

func NewAccountWrapper(cl Client) *AccountWrapper {
	return &AccountWrapper{
		cl: cl,
	}
}

func (w *AccountWrapper) GetAccountBySshKey(ctx context.Context, sshPubKey []byte) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/wrapper/GetAccountBySshKey")
	defer span.Finish()

	req := &pb.GetAccountBySshKeyRequest{
		SshPubKey: sshPubKey,
	}

	resp, err := w.cl.GetAccountBySshKey(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.GetAccountId(), nil
}
