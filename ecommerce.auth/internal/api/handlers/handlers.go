package handlers

import (
	"ecommerce.auth/internal/clients"
	"go.uber.org/zap"
)

type Handler struct {
	Auth   *AuthHandler
}

func NewHandler(logger *zap.Logger, userClient clients.UserClient) *Handler {
	return &Handler{
		Auth: NewAuthHandler(logger, userClient),
	}
}