package handlers

import (
	"go.uber.org/zap"
)

type Handler struct {
	Products *ProductsHandler
}

func NewHandler(logger *zap.Logger, s Service) *Handler {
	return &Handler{
		Products: NewProductsHandler(logger, s),
	}
}