package repository

import (
	"context"

	"ecommerce.users.manager/internal/models"
	"go.uber.org/zap"
)

func (r *userRepository) GetByUsernameOrEmail(ctx context.Context, logger *zap.Logger, identifier string) (*models.User, error) {
	var user models.User
	
	err := r.db.NewSelect().
		Model(&user).
		Where("username = ? OR email = ?", identifier, identifier).
		Scan(ctx)
	
	if err != nil {
		logger.Error("Failed to get user by username or email", zap.Error(err))
		return nil, err
	}
	
	return &user, nil
}