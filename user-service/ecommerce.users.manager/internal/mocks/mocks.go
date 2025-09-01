package mocks

import (
	"context"
	"ecommerce.users.manager/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type MockRepository struct {
	mock.Mock
}

// ---===[ genericRepository ]===---

func (m *MockRepository) SearchGeneric(ctx context.Context, logger *zap.Logger, model any, limit int, offset int, queryFunc func(bun.QueryBuilder) bun.QueryBuilder) (int, error) {
	args := m.Called(ctx, logger, model, limit, offset, queryFunc)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) CreateGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	args := m.Called(ctx, logger, model)
	return args.Error(0)
}

func (m *MockRepository) GetGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	args := m.Called(ctx, logger, model)
	return args.Error(0)
}

func (m *MockRepository) DeleteGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	args := m.Called(ctx, logger, model)
	return args.Error(0)
}

func (m *MockRepository) UpdateGeneric(ctx context.Context, logger *zap.Logger, model any) (int64, error) {
	args := m.Called(ctx, logger, model)
	return args.Get(0).(int64), args.Error(1)
}

// ---===[ usersRepository ]===---

func (m *MockRepository) GetByUsernameOrEmail(ctx context.Context, logger *zap.Logger, identifier string) (*models.User, error) {
	args := m.Called(ctx, logger, identifier)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockService struct {
	mock.Mock
}

// ---===[ usersService ]===---

func (m *MockService) CreateUser(ctx context.Context, logger *zap.Logger, req *models.CreateUserRequest) (*models.User, error) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockService) GetUser(ctx context.Context, logger *zap.Logger, req *models.GetUserRequest) (*models.User, error) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockService) UpdateUser(ctx context.Context, logger *zap.Logger, req *models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockService) DeleteUser(ctx context.Context, logger *zap.Logger, req *models.DeleteUserRequest) (bool, error) {
	args := m.Called(ctx, logger, req)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) AuthenticateUser(ctx context.Context, logger *zap.Logger, req *models.AuthRequest) (*models.AuthResponse, error) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.AuthResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockService) ValidateUser(ctx context.Context, logger *zap.Logger, req *models.ValidateUserRequest) (*models.ValidateUserResponse, error) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.ValidateUserResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
