package service

import (
	"context"
	"time"

	"ecommerce.products.manager/internal/errors"
	"ecommerce.products.manager/internal/models"
	"ecommerce.products.manager/internal/utils"

	"go.uber.org/zap"
)

func (s *ProductsService) CreateProduct(ctx context.Context, logger *zap.Logger, request *models.CreateProductRequest) (*models.Product, *errors.ErrorAPI) {
	logger.Info("Creating Product", zap.Any("request", request))

	if err := utils.ValidateProductFields(logger, &request.Name, &request.Price, &request.Stock, request.Description, true); err != nil {
		return nil, err
	}

	product := &models.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := s.productRepo.CreateGeneric(ctx, logger, product); err != nil {
		logger.Error("Failed to create product", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	logger.Info("Product created successfully", zap.String("id", product.ID.String()))
	return product, nil
}

func (s *ProductsService) GetProduct(ctx context.Context, logger *zap.Logger, request *models.GetProductRequest) (*models.Product, *errors.ErrorAPI) {
	logger.Info("Getting Product", zap.String("id", request.ID.String()))

	product := &models.Product{ID: request.ID}
	if err := s.productRepo.GetGeneric(ctx, logger, product); err != nil {
		logger.Error("Failed to get product", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	logger.Info("Product retrieved successfully", zap.Any("product", product))
	return product, nil
}

func (s *ProductsService) UpdateProduct(ctx context.Context, logger *zap.Logger, request *models.UpdateProductRequest) (*models.Product, *errors.ErrorAPI) {
	logger.Info("Updating Product", zap.Any("request", request))

	if err := utils.ValidateProductFields(logger, request.Name, request.Price, request.Stock, request.Description, true); err != nil {
		return nil, err
	}

	now := time.Now()
	request.UpdatedAt = now

	rowsAffected, err := s.productRepo.UpdateGeneric(ctx, logger, request)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	if rowsAffected == 0 {
		logger.Warn("No rows were updated, product may not exist", zap.Any("request", request))
		return nil, errors.NewErrorAPI(404, errors.ErrNotFound, "Product not found")
	}

	product := &models.Product{ID: request.ID}
	if err := s.productRepo.GetGeneric(ctx, logger, product); err != nil {
		logger.Error("Failed to retrieve updated product", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	logger.Info("Product updated successfully", zap.Any("product", product))
	return product, nil
}

func (s *ProductsService) DeleteProduct(ctx context.Context, logger *zap.Logger, request *models.DeleteProductRequest) (*models.DeleteProductResponse, *errors.ErrorAPI) {
	logger.Info("Deleting Product", zap.String("id", request.ID.String()))

	product := &models.Product{ID: request.ID}

	err := s.productRepo.DeleteGeneric(ctx, logger, product)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	logger.Info("Product deleted successfully", zap.String("id", request.ID.String()))
	return &models.DeleteProductResponse{Success: true}, nil
}

func (s *ProductsService) SearchProducts(ctx context.Context, logger *zap.Logger, request *models.SearchProductsRequest) ([]models.Product, int, *errors.ErrorAPI) {
	logger.Info("Searching Products", zap.Any("request", request))

	products, total, err := s.productRepo.SearchProducts(ctx, logger, request)
	if err != nil {
		logger.Error("Failed to search products", zap.Error(err))
		return nil, 0, errors.ValidateError(err)
	}

	logger.Info("Products search completed", zap.Int("total", total))
	return products, total, nil
}

func (s *ProductsService) UpdateStock(ctx context.Context, logger *zap.Logger, request *models.UpdateStockRequest) (*models.Product, *errors.ErrorAPI) {
	logger.Info("Updating Product Stock", zap.Any("request", request))

	if err := utils.ValidateProductFields(logger, nil, nil, &request.Stock, nil, false); err != nil {
		return nil, err
	}

	if err := s.productRepo.UpdateStock(ctx, logger, request); err != nil {
		logger.Error("Failed to update stock", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	product := &models.Product{ID: request.ID}
	if err := s.productRepo.GetGeneric(ctx, logger, product); err != nil {
		logger.Error("Failed to retrieve product after stock update", zap.Error(err))
		return nil, errors.ValidateError(err)
	}

	logger.Info("Stock updated successfully", zap.Any("product", product))
	return product, nil
}