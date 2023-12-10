package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}
