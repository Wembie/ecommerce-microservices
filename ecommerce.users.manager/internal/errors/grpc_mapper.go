package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPCError(err error) error {
	switch e := err.(type) {
	case Error:
		switch e.Code {
		case ErrCodeInvalidArgument:
			return status.Error(codes.InvalidArgument, e.Message)
		case ErrCodeNotFound:
			return status.Error(codes.NotFound, e.Message)
		case ErrCodeInternal:
			return status.Error(codes.Internal, e.Message)
		default:
			return status.Error(codes.Internal, "unexpected domain error: "+e.Message)
		}
	default:
		return status.Error(codes.Internal, err.Error())
	}
}