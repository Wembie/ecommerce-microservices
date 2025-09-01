package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"ecommerce.users.manager/internal/utils"
)

func TestStringPtr(t *testing.T) {
	input := "hello"
	ptr := utils.StringPtr(input)

	assert.NotNil(t, ptr, "Pointer should not be nil")
	assert.Equal(t, input, *ptr, "Pointer value should match the original string")
}

func TestTimePtr(t *testing.T) {
	now := time.Now()
	ptr := utils.TimePtr(now)

	assert.NotNil(t, ptr, "Pointer should not be nil")
	assert.Equal(t, now, *ptr, "Pointer value should match the original time")
}

func TestBoolPtr(t *testing.T) {
	input := true
	ptr := utils.BoolPtr(input)

	assert.NotNil(t, ptr, "Pointer should not be nil")
	assert.Equal(t, input, *ptr, "Pointer value should match the original boolean")
}
