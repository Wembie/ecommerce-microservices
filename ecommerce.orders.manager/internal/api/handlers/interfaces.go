package handlers

import (
	"context"

	"ecommerce.orders.manager/internal/errors"
	"ecommerce.orders.manager/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	OrdersService
}

type OrdersService interface {
	CreateOrder(context.Context, *zap.Logger, *models.CreateOrderRequest) (*models.OrderResponse, *errors.ErrorAPI)
	GetOrder(context.Context, *zap.Logger, *models.GetOrderRequest) (*models.OrderResponse, *errors.ErrorAPI)
	GetOrdersByUser(context.Context, *zap.Logger, *models.GetOrdersByUserRequest) ([]*models.Order, int, *errors.ErrorAPI)
	UpdateOrderStatus(context.Context, *zap.Logger, *models.UpdateOrderStatusRequest) (*models.Order, *errors.ErrorAPI)
	GetOrderItems(context.Context, *zap.Logger, uuid.UUID) ([]*models.OrderItem, *errors.ErrorAPI)
}