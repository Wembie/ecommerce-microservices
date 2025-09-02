package repository

import (
	"ecommerce.users.manager/internal/config"
	"ecommerce.users.manager/internal/service"

	"github.com/uptrace/bun"
)

type userRepository struct {
	db     *bun.DB
	config *config.Config
}

var _ service.UserRepository = (*userRepository)(nil)

func NewRepository(db *bun.DB, cfg *config.Config) *userRepository {
	return &userRepository{
		db:     db,
		config: cfg,
	}
}