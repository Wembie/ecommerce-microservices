package transport

import (
	"context"

	"ecommerce.users.manager/internal/models"
	"go.uber.org/zap"
)

type Service interface {
	usersService
}

type usersService interface {
	CreateUser(context.Context, *zap.Logger, *models.CreateUserRequest) (*models.User, error)
	GetUser(context.Context, *zap.Logger, *models.GetUserRequest) (*models.User, error)
	UpdateUser(context.Context, *zap.Logger, *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(context.Context, *zap.Logger, *models.DeleteUserRequest) (bool, error)
	AuthenticateUser(context.Context, *zap.Logger, *models.AuthRequest) (*models.AuthResponse, error)
	ValidateUser(context.Context, *zap.Logger, *models.ValidateUserRequest) (*models.ValidateUserResponse, error)
}