package repository

import (
	"ecommerce.products.manager/internal/config"
	"ecommerce.products.manager/internal/service"

	"github.com/uptrace/bun"
)

type productRepository struct {
	db     *bun.DB
	config *config.Config
}

var _ service.ProductRepository = (*productRepository)(nil)

func NewRepository(db *bun.DB, cfg *config.Config) *productRepository {
	return &productRepository{
		db:     db,
		config: cfg,
	}
}