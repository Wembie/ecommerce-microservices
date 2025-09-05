package service

import (
	"context"
	"fmt"
	"time"

	"ecommerce.orders.manager/internal/errors"
	"ecommerce.orders.manager/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *OrdersService) CreateOrder(
	ctx context.Context,
	logger *zap.Logger,
	request *models.CreateOrderRequest,
) (*models.OrderResponse, *errors.ErrorAPI) {
	logger.Info("Service: CreateOrder")

	if len(request.Items) == 0 {
		return nil, errors.ValidateError(errors.ErrInvalidInput, "Order must have at least one item")
	}

	var total float64
	var orderItems []*models.OrderItem
	orderID := uuid.New()

	for _, item := range request.Items {
		if item.Quantity <= 0 {
			return nil, errors.ValidateError(errors.ErrInvalidInput, "Item quantity must be greater than 0")
		}

		product, err := s.productClient.GetProduct(ctx, item.ProductID)
		if err != nil {
			logger.Error("Failed to get product", zap.Error(err))
			return nil, errors.ValidateError(errors.ErrNotFound, "Product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.ValidateError(errors.ErrInsufficientStock, fmt.Sprintf("Not enough stock for product %s", product.Name))
		}

		orderItem := &models.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}

		orderItems = append(orderItems, orderItem)
		total += product.Price * float64(item.Quantity)

		newStock := product.Stock - item.Quantity
		if err := s.productClient.UpdateStock(ctx, item.ProductID, newStock); err != nil {
			logger.Error("Failed to update product stock", zap.Error(err))
			return nil, errors.ValidateError(errors.ErrInternalServer, "Failed to update product stock")
		}
	}

	order := &models.Order{
		ID:     orderID,
		UserID: request.UserID,
		Status: models.OrderStatuses.Pending,
		Total:  total,
	}

	if err := s.orderRepo.CreateOrderWithItems(ctx, logger, order, orderItems); err != nil {
		logger.Error("Failed to create order", zap.Error(err))
		return nil, errors.ValidateError(errors.ErrInternalServer, "Failed to create order")
	}

	return &models.OrderResponse{
		Order: order,
		Items: orderItems,
	}, nil
}

func (s *OrdersService) GetOrder(
	ctx context.Context,
	logger *zap.Logger,
	request *models.GetOrderRequest,
) (*models.OrderResponse, *errors.ErrorAPI) {
	logger.Info("Service: GetOrder", zap.String("order_id", request.ID.String()))

	order, items, err := s.orderRepo.GetOrderWithItems(ctx, logger, request.ID)
	if err != nil {
		logger.Error("Failed to get order", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	return &models.OrderResponse{
		Order: order,
		Items: items,
	}, nil
}

func (s *OrdersService) GetOrdersByUser(
	ctx context.Context,
	logger *zap.Logger,
	request *models.GetOrdersByUserRequest,
) ([]*models.Order, int, *errors.ErrorAPI) {
	logger.Info("Service: GetOrdersByUser", zap.String("user_id", request.UserID.String()))

	orders, total, err := s.orderRepo.GetOrdersByUser(ctx, logger, request)
	if err != nil {
		logger.Error("Failed to get orders by user", zap.Error(err))
		return nil, 0, errors.ValidateError(err)
	}

	return orders, total, nil
}

func (s *OrdersService) UpdateOrderStatus(
	ctx context.Context,
	logger *zap.Logger,
	request *models.UpdateOrderStatusRequest,
) (*models.Order, *errors.ErrorAPI) {
	logger.Info("Service: UpdateOrderStatus", zap.String("order_id", request.ID.String()))

	validStatuses := []string{
		models.OrderStatuses.Pending,
		models.OrderStatuses.Confirmed,
		models.OrderStatuses.Shipped,
		models.OrderStatuses.Delivered,
		models.OrderStatuses.Cancelled,
	}

	isValid := false
	for _, status := range validStatuses {
		if request.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.ValidateError(errors.ErrInvalidInput, "Invalid order status")
	}

	request.UpdatedAt = time.Now()

	if err := s.orderRepo.UpdateOrderStatus(ctx, logger, request); err != nil {
		logger.Error("Failed to update order status", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	var order models.Order
	order.ID = request.ID
	if err := s.orderRepo.GetGeneric(ctx, logger, &order); err != nil {
		logger.Error("Failed to get updated order", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	return &order, nil
}

func (s *OrdersService) GetOrderItems(
	ctx context.Context,
	logger *zap.Logger,
	orderID uuid.UUID,
) ([]*models.OrderItem, *errors.ErrorAPI) {
	logger.Info("Service: GetOrderItems", zap.String("order_id", orderID.String()))

	items, err := s.orderRepo.GetOrderItems(ctx, logger, orderID)
	if err != nil {
		logger.Error("Failed to get order items", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	return items, nil
}