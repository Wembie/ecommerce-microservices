package mocks_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ecommerce.users.manager/internal/mocks"
)

func TestNewTestSetup(t *testing.T) {
	ts := mocks.NewTestSetup()

	assert.NotNil(t, ts)
	assert.NotNil(t, ts.Logger)
	assert.NotNil(t, ts.Ctx)
	assert.NotNil(t, ts.Config)
	assert.Equal(t, "test_jwt_secret", ts.Config.JWTSecret)
}

func TestNewMockRepository(t *testing.T) {
	ts := mocks.NewTestSetup()
	mockRepo := ts.NewMockRepository()

	assert.NotNil(t, mockRepo)
}

func TestNewMockService(t *testing.T) {
	ts := mocks.NewTestSetup()
	mockService := ts.NewMockService()

	assert.NotNil(t, mockService)
}
