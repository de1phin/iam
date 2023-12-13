package service

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorNotFound(entity string, id string) error {
	return status.Error(codes.NotFound, fmt.Sprintf("%s with ID %s not found", entity, id))
}

func ErrorAccountNotFound(id string) error {
	return ErrorNotFound("account", id)
}

func ErrorInternal(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func ErrorAlreadyExists(obj string) error {
	return status.Error(codes.AlreadyExists, obj+" already exists")
}
