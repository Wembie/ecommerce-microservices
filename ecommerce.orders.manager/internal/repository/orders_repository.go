package repository

import (
	"context"
	"time"

	"ecommerce.orders.manager/internal/models"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

func (r *orderRepository) CreateOrderWithItems(
	ctx context.Context,
	logger *zap.Logger,
	order *models.Order,
	items []*models.OrderItem,
) error {
	logger.Info("Repository: CreateOrderWithItems")

	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(order).
			Exec(ctx)
		if err != nil {
			logger.Error("Failed to create order", zap.Error(err))
			return err
		}

		if len(items) > 0 {
			_, err = tx.NewInsert().
				Model(&items).
				Exec(ctx)
			if err != nil {
				logger.Error("Failed to create order items", zap.Error(err))
				return err
			}
		}

		return nil
	})
}

func (r *orderRepository) GetOrderWithItems(
	ctx context.Context,
	logger *zap.Logger,
	orderID uuid.UUID,
) (*models.Order, []*models.OrderItem, error) {
	logger.Info("Repository: GetOrderWithItems", zap.String("order_id", orderID.String()))

	var order models.Order
	order.ID = orderID

	err := r.GetGeneric(ctx, logger, &order)
	if err != nil {
		return nil, nil, err
	}

	var items []*models.OrderItem
	err = r.db.NewSelect().
		Model(&items).
		Where("order_id = ?", orderID).
		Scan(ctx)
	if err != nil {
		logger.Error("Failed to get order items", zap.Error(err))
		return nil, nil, err
	}

	return &order, items, nil
}

func (r *orderRepository) GetOrdersByUser(
	ctx context.Context,
	logger *zap.Logger,
	request *models.GetOrdersByUserRequest,
) ([]*models.Order, int, error) {
	logger.Info("Repository: GetOrdersByUser")

	var orders []*models.Order

	queryModifier := func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where("user_id = ?", request.UserID).
			Order("created_at DESC")
	}

	total, err := r.SearchGeneric(ctx, logger, &orders, request.Page, request.Size, queryModifier)
	if err != nil {
		logger.Error("Failed to get orders by user", zap.Error(err))
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) UpdateOrderStatus(
	ctx context.Context,
	logger *zap.Logger,
	request *models.UpdateOrderStatusRequest,
) error {
	logger.Info("Repository: UpdateOrderStatus", zap.String("id", request.ID.String()), zap.String("status", request.Status))

	now := time.Now()
	_, err := r.db.NewUpdate().
		Table("order_service.orders").
		Set("status = ?", request.Status).
		Set("updated_at = ?", now).
		Where("id = ?", request.ID).
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to update order status", zap.Error(err))
		return err
	}

	return nil
}

func (r *orderRepository) GetOrderItems(
	ctx context.Context,
	logger *zap.Logger,
	orderID uuid.UUID,
) ([]*models.OrderItem, error) {
	logger.Info("Repository: GetOrderItems", zap.String("order_id", orderID.String()))

	var items []*models.OrderItem
	err := r.db.NewSelect().
		Model(&items).
		Where("order_id = ?", orderID).
		Scan(ctx)

	if err != nil {
		logger.Error("Failed to get order items", zap.Error(err))
		return nil, err
	}

	return items, nil
}