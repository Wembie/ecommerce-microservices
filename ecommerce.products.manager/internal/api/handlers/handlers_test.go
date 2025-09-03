package handlers_test

import (
	"testing"

	"ecommerce.products.manager/internal/mocks"
	"ecommerce.products.manager/internal/api/handlers"

	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	setup := mocks.NewTestSetup()
	mockSvc := setup.NewMockService()

	handler := handlers.NewHandler(setup.Logger, mockSvc)

	assert.NotNil(t, handler)
	assert.IsType(t, &handlers.Handler{}, handler)
}
