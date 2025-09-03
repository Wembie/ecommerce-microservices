package service_test

import (
	"testing"
	
	"ecommerce.users.manager/internal/service"
	"ecommerce.users.manager/internal/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	setup := mocks.NewTestSetup()
	mockRepo := setup.NewMockRepository()

	svc := service.NewService(setup.Logger, mockRepo, setup.Config)

	assert.NotNil(t, svc)
	assert.IsType(t, &service.UsersService{}, svc)
}
