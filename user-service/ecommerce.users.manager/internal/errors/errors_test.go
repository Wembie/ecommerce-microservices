package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ecommerce.users.manager/internal/errors"
)

func TestErrorImplementsErrorInterface(t *testing.T) {
	err := errors.Error{Code: errors.ErrCodeInvalidArgument, Message: "invalid arg"}
	var stdErr error = err
	assert.Equal(t, "invalid arg", stdErr.Error())
}

func TestNewReturnsError(t *testing.T) {
	err := errors.New(errors.ErrCodeNotFound, "not found")
	assert.Error(t, err)

	e, ok := err.(errors.Error)
	assert.True(t, ok)
	assert.Equal(t, errors.ErrCodeNotFound, e.Code)
	assert.Equal(t, "not found", e.Message)
}
