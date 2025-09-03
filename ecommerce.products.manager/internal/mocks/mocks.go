package mocks

import (
	"context"

	"ecommerce.products.manager/internal/errors"
	"ecommerce.products.manager/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type MockRepository struct {
	mock.Mock
}

// ---===[ genericRepository ]===---

func (m *MockRepository) SearchGeneric(
	ctx context.Context,
	logger *zap.Logger,
	model any,
	limit int,
	offset int,
	queryFunc func(*bun.SelectQuery) *bun.SelectQuery,
) (int, error) {
	args := m.Called(ctx, logger, model, limit, offset, queryFunc)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) CreateGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	args := m.Called(ctx, logger, model)
	return args.Error(0)
}

func (m *MockRepository) GetGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	args := m.Called(ctx, logger, model)
	return args.Error(0)
}

func (m *MockRepository) DeleteGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	args := m.Called(ctx, logger, model)
	return args.Error(0)
}

func (m *MockRepository) UpdateGeneric(ctx context.Context, logger *zap.Logger, model any) (int64, error) {
	args := m.Called(ctx, logger, model)
	return args.Get(0).(int64), args.Error(1)
}

// ---===[ productsRepository ]===---

func (m *MockRepository) SearchProducts(ctx context.Context, logger *zap.Logger, req *models.SearchProductsRequest) ([]models.Product, int, error) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Product), args.Int(1), args.Error(2)
	}
	return nil, args.Int(1), args.Error(2)
}

func (m *MockRepository) UpdateStock(ctx context.Context, logger *zap.Logger, req *models.UpdateStockRequest) error {
	args := m.Called(ctx, logger, req)
	return args.Error(0)
}


type MockService struct {
	mock.Mock
}

// ---===[ ProductsService ]===---

func (m *MockService) CreateProduct(ctx context.Context, logger *zap.Logger, req *models.CreateProductRequest) (*models.Product, *errors.ErrorAPI) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Get(1).(*errors.ErrorAPI)
	}
	return nil, args.Get(1).(*errors.ErrorAPI)
}

func (m *MockService) GetProduct(ctx context.Context, logger *zap.Logger, req *models.GetProductRequest) (*models.Product, *errors.ErrorAPI) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Get(1).(*errors.ErrorAPI)
	}
	return nil, args.Get(1).(*errors.ErrorAPI)
}

func (m *MockService) UpdateProduct(ctx context.Context, logger *zap.Logger, req *models.UpdateProductRequest) (*models.Product, *errors.ErrorAPI) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Get(1).(*errors.ErrorAPI)
	}
	return nil, args.Get(1).(*errors.ErrorAPI)
}

func (m *MockService) DeleteProduct(ctx context.Context, logger *zap.Logger, req *models.DeleteProductRequest) (*models.DeleteProductResponse, *errors.ErrorAPI) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.DeleteProductResponse), args.Get(1).(*errors.ErrorAPI)
	}
	return nil, args.Get(1).(*errors.ErrorAPI)
}

func (m *MockService) SearchProducts(ctx context.Context, logger *zap.Logger, req *models.SearchProductsRequest) ([]models.Product, int, *errors.ErrorAPI) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Product), args.Int(1), args.Get(2).(*errors.ErrorAPI)
	}
	return nil, args.Int(1), args.Get(2).(*errors.ErrorAPI)
}

func (m *MockService) UpdateStock(ctx context.Context, logger *zap.Logger, req *models.UpdateStockRequest) (*models.Product, *errors.ErrorAPI) {
	args := m.Called(ctx, logger, req)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Product), args.Get(1).(*errors.ErrorAPI)
	}
	return nil, args.Get(1).(*errors.ErrorAPI)
}
