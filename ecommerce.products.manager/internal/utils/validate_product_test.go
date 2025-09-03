package utils_test

import (
	"net/http"
	"testing"

	"ecommerce.products.manager/internal/errors"
	"ecommerce.products.manager/internal/utils"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestValidateProductFields(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Name required but nil", func(t *testing.T) {
		err := utils.ValidateProductFields(logger, nil, nil, nil, nil, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrInvalidInput, err.Error)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	})

	t.Run("Name empty", func(t *testing.T) {
		name := ""
		err := utils.ValidateProductFields(logger, &name, nil, nil, nil, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrInvalidInput, err.Error)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	})

	t.Run("Name too long", func(t *testing.T) {
		longName := make([]byte, 256)
		for i := range longName {
			longName[i] = 'a'
		}
		name := string(longName)
		err := utils.ValidateProductFields(logger, &name, nil, nil, nil, true)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrInvalidInput, err.Error)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	})

	t.Run("Invalid price", func(t *testing.T) {
		price := -10.0
		err := utils.ValidateProductFields(logger, nil, &price, nil, nil, false)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrInvalidInput, err.Error)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	})

	t.Run("Invalid stock", func(t *testing.T) {
		stock := -5
		err := utils.ValidateProductFields(logger, nil, nil, &stock, nil, false)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrInvalidInput, err.Error)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	})

	t.Run("Description too long", func(t *testing.T) {
		longDesc := make([]byte, 1001)
		for i := range longDesc {
			longDesc[i] = 'x'
		}
		desc := string(longDesc)
		err := utils.ValidateProductFields(logger, nil, nil, nil, &desc, false)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrInvalidInput, err.Error)
		assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	})

	t.Run("Valid case with all fields", func(t *testing.T) {
		name := "Valid Name"
		price := 100.0
		stock := 10
		desc := "A valid product description"
		err := utils.ValidateProductFields(logger, &name, &price, &stock, &desc, true)
		assert.Nil(t, err)
	})

	t.Run("Valid update without name", func(t *testing.T) {
		price := 50.0
		err := utils.ValidateProductFields(logger, nil, &price, nil, nil, false)
		assert.Nil(t, err)
	})
}