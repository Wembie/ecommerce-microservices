package mocks

import (
	"context"

	"ecommerce.users.manager/internal/config"
	"go.uber.org/zap"
)

type TestSetup struct {
	Logger *zap.Logger
	Ctx    context.Context
	Config *config.Config
}

func NewTestSetup() *TestSetup {
	return &TestSetup{
		Logger: zap.NewNop(),
		Ctx:    context.Background(),
		Config: &config.Config{
			JWTSecret: "test_jwt_secret",
		},
	}
}

func (ts *TestSetup) NewMockRepository() *MockRepository {
	return &MockRepository{}
}

func (ts *TestSetup) NewMockService() *MockService {
	return &MockService{}
}
