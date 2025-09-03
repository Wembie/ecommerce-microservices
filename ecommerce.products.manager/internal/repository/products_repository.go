package repository

import (
	"context"
	"fmt"
	"time"

	"ecommerce.products.manager/internal/models"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

func (r *productRepository) SearchProducts(
	ctx context.Context,
	logger *zap.Logger,
	request *models.SearchProductsRequest,
) ([]models.Product, int, error) {
	logger.Info("Repository: SearchProducts")

	var products []models.Product

	queryModifier := func(q *bun.SelectQuery) *bun.SelectQuery {
		if request.Name != nil && *request.Name != "" {
			q = q.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", *request.Name))
		}

		if request.Description != nil && *request.Description != "" {
			q = q.Where("description ILIKE ?", fmt.Sprintf("%%%s%%", *request.Description))
		}

		if request.Price != nil {
			q = q.Where("price = ?", *request.Price)
		}

		if request.Stock != nil {
			q = q.Where("stock >= ?", *request.Stock)
		}

		return q.Order("created_at DESC")
	}

	total, err := r.SearchGeneric(ctx, logger, &products, request.Page, request.Size, queryModifier)
	if err != nil {
		logger.Error("Failed to search products", zap.Error(err))
		return nil, 0, err
	}

	return products, total, nil
}


func (r *productRepository) UpdateStock(
	ctx context.Context,
	logger *zap.Logger,
	request *models.UpdateStockRequest,
) error {
	logger.Info("Repository: UpdateStock", zap.String("id", request.ID.String()), zap.Int("stock", request.Stock))

	now := time.Now()
	_, err := r.db.NewUpdate().
		Table("product_service.products").
		Set("stock = ?", request.Stock).
		Set("updated_at = ?", now).
		Where("id = ?", request.ID).
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to update stock", zap.Error(err))
		return err
	}

	return nil
}