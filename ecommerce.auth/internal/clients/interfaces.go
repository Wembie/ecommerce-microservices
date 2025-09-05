package clients

import (
	"context"
	"ecommerce.auth/internal/models"
)

type UserClient interface {
	AuthenticateUser(ctx context.Context, req *models.AuthRequest) (*models.AuthResponse, error)
	ValidateUser(ctx context.Context, token string) (*models.ValidateUserResponse, error)
	Close() error
}