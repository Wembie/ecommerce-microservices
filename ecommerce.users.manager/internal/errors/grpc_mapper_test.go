package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ecommerce.users.manager/internal/errors"
)

func TestToGRPCError_KnownCodes(t *testing.T) {
	tests := []struct {
		name     string
		code     errors.CodeError
		message  string
		expected codes.Code
	}{
		{"InvalidArgument", errors.ErrCodeInvalidArgument, "bad input", codes.InvalidArgument},
		{"NotFound", errors.ErrCodeNotFound, "missing", codes.NotFound},
		{"Internal", errors.ErrCodeInternal, "oops", codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Error{Code: tt.code, Message: tt.message}
			grpcErr := errors.ToGRPCError(err)
			st, ok := status.FromError(grpcErr)
			assert.True(t, ok)
			assert.Equal(t, tt.expected, st.Code())
			assert.Equal(t, tt.message, st.Message())
		})
	}
}

func TestToGRPCError_UnexpectedDomainError(t *testing.T) {
	err := errors.Error{Code: 999, Message: "weird"}
	grpcErr := errors.ToGRPCError(err)
	st, ok := status.FromError(grpcErr)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "unexpected domain error")
}

func TestToGRPCError_StandardError(t *testing.T) {
	_ = errors.New(errors.ErrCodeInternal, "something went wrong")
	stdErr := assert.AnError

	grpcErr := errors.ToGRPCError(stdErr)
	st, ok := status.FromError(grpcErr)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, stdErr.Error(), st.Message())
}
