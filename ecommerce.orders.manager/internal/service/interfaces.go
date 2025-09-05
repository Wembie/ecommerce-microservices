package service

import (
	"context"

	"ecommerce.orders.manager/internal/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type OrderRepository interface {
	genericRepository
	ordersRepository
}

type genericRepository interface {
	SearchGeneric(context.Context, *zap.Logger, any, int, int, func(*bun.SelectQuery) *bun.SelectQuery) (int, error)
	CreateGeneric(context.Context, *zap.Logger, any) error
	GetGeneric(context.Context, *zap.Logger, any) error
	DeleteGeneric(context.Context, *zap.Logger, any) error
	UpdateGeneric(context.Context, *zap.Logger, any) (int64, error)
}

type ordersRepository interface {
	CreateOrderWithItems(context.Context, *zap.Logger, *models.Order, []*models.OrderItem) error
	GetOrderWithItems(context.Context, *zap.Logger, uuid.UUID) (*models.Order, []*models.OrderItem, error)
	GetOrdersByUser(context.Context, *zap.Logger, *models.GetOrdersByUserRequest) ([]*models.Order, int, error)
	UpdateOrderStatus(context.Context, *zap.Logger, *models.UpdateOrderStatusRequest) error
	GetOrderItems(context.Context, *zap.Logger, uuid.UUID) ([]*models.OrderItem, error)
}