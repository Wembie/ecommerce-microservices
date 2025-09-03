package handlers

import (
	"context"

	"ecommerce.products.manager/internal/errors"
	"ecommerce.products.manager/internal/models"
	"go.uber.org/zap"
)

type Service interface {
	ProductsService
}

type ProductsService interface {
	CreateProduct(context.Context, *zap.Logger, *models.CreateProductRequest) (*models.Product, *errors.ErrorAPI)
	GetProduct(context.Context, *zap.Logger, *models.GetProductRequest) (*models.Product, *errors.ErrorAPI)
	UpdateProduct(context.Context, *zap.Logger, *models.UpdateProductRequest) (*models.Product, *errors.ErrorAPI)
	DeleteProduct(context.Context, *zap.Logger, *models.DeleteProductRequest) (*models.DeleteProductResponse, *errors.ErrorAPI)
	SearchProducts(context.Context, *zap.Logger, *models.SearchProductsRequest) ([]models.Product, int, *errors.ErrorAPI)
	UpdateStock(context.Context, *zap.Logger, *models.UpdateStockRequest) (*models.Product, *errors.ErrorAPI)
}