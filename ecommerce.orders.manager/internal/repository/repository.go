package repository

import (
	"ecommerce.orders.manager/internal/config"
	"ecommerce.orders.manager/internal/service"

	"github.com/uptrace/bun"
)

type orderRepository struct {
	db     *bun.DB
	config *config.Config
}

var _ service.OrderRepository = (*orderRepository)(nil)

func NewRepository(db *bun.DB, cfg *config.Config) *orderRepository {
	return &orderRepository{
		db:     db,
		config: cfg,
	}
}