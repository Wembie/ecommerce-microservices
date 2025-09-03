package service

import (
	"context"

	"ecommerce.products.manager/internal/models"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type ProductRepository interface {
	genericRepository
	productsRepository
}

type genericRepository interface {
	SearchGeneric(context.Context, *zap.Logger, any, int, int, func(*bun.SelectQuery) *bun.SelectQuery) (int, error)
	CreateGeneric(context.Context, *zap.Logger, any) error
	GetGeneric(context.Context, *zap.Logger, any) error
	DeleteGeneric(context.Context, *zap.Logger, any) error
	UpdateGeneric(context.Context, *zap.Logger, any) (int64, error)
}

type productsRepository interface {
	SearchProducts(context.Context, *zap.Logger, *models.SearchProductsRequest) ([]models.Product, int, error)
	UpdateStock(context.Context, *zap.Logger, *models.UpdateStockRequest) error
}