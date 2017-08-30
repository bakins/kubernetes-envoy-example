package util

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewNotFoundError creates a new NotFound grpc error
func NewNotFoundError(kind, name string) error {
	return grpc.Errorf(codes.NotFound, "not found: %s/%s", kind, name)
}

// IsNotFoundError determines if the error is a NotFound grpc error
func IsNotFoundError(err error) bool {
	return isErrorCode(err, codes.NotFound)
}

// NewAlreadyExistsError creates a new AlreadyExists grpc error
func NewAlreadyExistsError(kind, name string) error {
	return grpc.Errorf(codes.AlreadyExists, "already exists: %s/%s", kind, name)
}

// IsAlreadyExistsError determines if the error is a AlreadyExists grpc error
func IsAlreadyExistsError(err error) bool {
	return isErrorCode(err, codes.AlreadyExists)
}

// NewInvalidArgumentError creates a new InvalidArgument grpc error
func NewInvalidArgumentError(kind, name string) error {
	return grpc.Errorf(codes.InvalidArgument, "invalid argument: %s/%s", kind, name)
}

// IsInvalidArgumentError determines if the error is a InvalidArgument grpc error
func IsInvalidArgumentError(err error) bool {
	return isErrorCode(err, codes.InvalidArgument)
}

func isErrorCode(err error, c codes.Code) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	return s.Code() == c
}
