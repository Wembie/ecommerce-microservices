package handlers

import (
	"ecommerce.orders.manager/internal/clients"
	"go.uber.org/zap"
)

type Handler struct {
	Orders *OrdersHandler
	Auth   *AuthHandler
}

func NewHandler(logger *zap.Logger, s Service, userClient clients.UserClient) *Handler {
	return &Handler{
		Orders: NewOrdersHandler(logger, s),
		Auth:   NewAuthHandler(logger, userClient),
	}
}