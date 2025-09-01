package transport_test

import (
	"ecommerce.users.manager/internal/mocks"
	"ecommerce.users.manager/internal/transport"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	setup := mocks.NewTestSetup()
	mockSvc := setup.NewMockService()

	handler := transport.NewHandler(setup.Logger, mockSvc, setup.Config)

	assert.NotNil(t, handler)
	assert.IsType(t, &transport.PingHandler{}, handler)
}