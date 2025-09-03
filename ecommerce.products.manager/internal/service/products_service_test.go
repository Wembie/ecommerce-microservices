package service_test

import (
	"testing"
	"time"

	"ecommerce.products.manager/internal/errors"
	"ecommerce.products.manager/internal/mocks"
	"ecommerce.products.manager/internal/models"
	"ecommerce.products.manager/internal/service"
	"ecommerce.products.manager/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductsService_CreateProduct(t *testing.T) {
	setup := mocks.NewTestSetup()

	tests := []struct {
		name        string
		request     *models.CreateProductRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
	}{
		{
			name: "successful product creation",
			request: &models.CreateProductRequest{
				Name:        "Test Product",
				Description: utils.StringPtr("Test Description"),
				Price:       99.99,
				Stock:       10,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("CreateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.Product)
					product.ID = uuid.New()
					product.CreatedAt = time.Now()
				})
			},
			expectError: false,
		},
		{
			name: "successful product creation without description",
			request: &models.CreateProductRequest{
				Name:  "Test Product",
				Price: 99.99,
				Stock: 10,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("CreateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.Product)
					product.ID = uuid.New()
					product.CreatedAt = time.Now()
				})
			},
			expectError: false,
		},
		{
			name: "repository error should return error",
			request: &models.CreateProductRequest{
				Name:        "Test Product",
				Description: utils.StringPtr("Test Description"),
				Price:       99.99,
				Stock:       10,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("CreateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(errors.ErrInternalServer)
			},
			expectError: true,
		},
		{
			name: "validation error - empty name",
			request: &models.CreateProductRequest{
				Name:  "",
				Price: 99.99,
				Stock: 10,
			},
			setupMocks:  func(repo *mocks.MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - negative price",
			request: &models.CreateProductRequest{
				Name:  "Test Product",
				Price: -10.50,
				Stock: 10,
			},
			setupMocks:  func(repo *mocks.MockRepository) {},
			expectError: true,
		},
		{
			name: "validation error - negative stock",
			request: &models.CreateProductRequest{
				Name:  "Test Product",
				Price: 99.99,
				Stock: -5,
			},
			setupMocks:  func(repo *mocks.MockRepository) {},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.CreateProduct(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, tt.request.Description, result.Description)
				assert.Equal(t, tt.request.Price, result.Price)
				assert.Equal(t, tt.request.Stock, result.Stock)
				assert.NotEmpty(t, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductsService_GetProduct(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()

	tests := []struct {
		name        string
		request     *models.GetProductRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
	}{
		{
			name: "successful get product",
			request: &models.GetProductRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.Product)
					product.ID = testID
					product.Name = "Test Product"
					product.Description = utils.StringPtr("Test Description")
					product.Price = 99.99
					product.Stock = 10
					product.CreatedAt = time.Now()
					product.UpdatedAt = utils.TimePtr(time.Now())
				})
			},
			expectError: false,
		},
		{
			name: "product not found should return error",
			request: &models.GetProductRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(errors.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.GetProduct(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testID, result.ID)
				assert.Equal(t, "Test Product", result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductsService_UpdateProduct(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()
	newName := "Updated Product"
	newDescription := "Updated Description"
	newPrice := 149.99
	newStock := 15

	tests := []struct {
		name        string
		request     *models.UpdateProductRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
		expectNil   bool
	}{
		{
			name: "successful update with all fields",
			request: &models.UpdateProductRequest{
				ID:          testID,
				Name:        utils.StringPtr(newName),
				Description: utils.StringPtr(newDescription),
				Price:       utils.Float64Ptr(newPrice),
				Stock:       utils.IntPtr(newStock),
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateProductRequest")).Return(int64(1), nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.Product)
					product.ID = testID
					product.Name = newName
					product.Description = utils.StringPtr(newDescription)
					product.Price = newPrice
					product.Stock = newStock
					product.CreatedAt = time.Now()
					product.UpdatedAt = utils.TimePtr(time.Now())
				})
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "successful partial update",
			request: &models.UpdateProductRequest{
				ID:   testID,
				Name: utils.StringPtr(newName),
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateProductRequest")).Return(int64(1), nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.Product)
					product.ID = testID
					product.Name = newName
					product.Description = utils.StringPtr("Original Description")
					product.Price = 99.99
					product.Stock = 10
					product.CreatedAt = time.Now()
					product.UpdatedAt = utils.TimePtr(time.Now())
				})
			},
			expectError: false,
			expectNil:   false,
		},
		{
			name: "no rows affected should return error",
			request: &models.UpdateProductRequest{
				ID:   testID,
				Name: utils.StringPtr(newName),
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateProductRequest")).Return(int64(0), nil)
			},
			expectError: true,
			expectNil:   false,
		},
		{
			name: "update error should return error",
			request: &models.UpdateProductRequest{
				ID:   testID,
				Name: utils.StringPtr(newName),
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateProductRequest")).Return(int64(0), errors.ErrInternalServer)
			},
			expectError: true,
			expectNil:   false,
		},
		{
			name: "get after update error should return error",
			request: &models.UpdateProductRequest{
				ID:   testID,
				Name: utils.StringPtr(newName),
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateProductRequest")).Return(int64(1), nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(errors.ErrInternalServer)
			},
			expectError: true,
			expectNil:   false,
		},
		{
			name: "validation error - negative price",
			request: &models.UpdateProductRequest{
				ID:    testID,
				Price: utils.Float64Ptr(-10.50),
			},
			setupMocks:  func(repo *mocks.MockRepository) {},
			expectError: true,
			expectNil:   false,
		},
		{
			name: "validation error - negative stock",
			request: &models.UpdateProductRequest{
				ID:    testID,
				Stock: utils.IntPtr(-5),
			},
			setupMocks:  func(repo *mocks.MockRepository) {},
			expectError: true,
			expectNil:   false,
		},
		{
			name: "validation error - empty name",
			request: &models.UpdateProductRequest{
				ID:   testID,
				Name: utils.StringPtr(""),
			},
			setupMocks:  func(repo *mocks.MockRepository) {},
			expectError: true,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.UpdateProduct(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				if tt.expectNil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, testID, result.ID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductsService_DeleteProduct(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()

	tests := []struct {
		name        string
		request     *models.DeleteProductRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
	}{
		{
			name: "successful delete",
			request: &models.DeleteProductRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("DeleteGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "repository error should return error",
			request: &models.DeleteProductRequest{
				ID: testID,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("DeleteGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(errors.ErrInternalServer)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.DeleteProduct(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.True(t, result.Success)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductsService_SearchProducts(t *testing.T) {
	setup := mocks.NewTestSetup()

	tests := []struct {
		name         string
		request      *models.SearchProductsRequest
		setupMocks   func(*mocks.MockRepository)
		expectError  bool
		expectedLen  int
		expectedTotal int
	}{
		{
			name: "successful search with results",
			request: &models.SearchProductsRequest{
				Name: utils.StringPtr("Test"),
				Page: 0,
				Size: 10,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				products := []models.Product{
					{
						ID:          uuid.New(),
						Name:        "Test Product 1",
						Description: utils.StringPtr("Test Description 1"),
						Price:       99.99,
						Stock:       10,
						CreatedAt:   time.Now(),
						UpdatedAt:   utils.TimePtr(time.Now()),
					},
					{
						ID:          uuid.New(),
						Name:        "Test Product 2",
						Description: utils.StringPtr("Test Description 2"),
						Price:       149.99,
						Stock:       5,
						CreatedAt:   time.Now(),
						UpdatedAt:   utils.TimePtr(time.Now()),
					},
				}
				repo.On("SearchProducts", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.SearchProductsRequest")).Return(products, 2, nil)
			},
			expectError:   false,
			expectedLen:   2,
			expectedTotal: 2,
		},
		{
			name: "successful search with no results",
			request: &models.SearchProductsRequest{
				Name: utils.StringPtr("NonExistent"),
				Page: 0,
				Size: 10,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				products := []models.Product{}
				repo.On("SearchProducts", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.SearchProductsRequest")).Return(products, 0, nil)
			},
			expectError:   false,
			expectedLen:   0,
			expectedTotal: 0,
		},
		{
			name: "successful search with multiple filters",
			request: &models.SearchProductsRequest{
				Name:        utils.StringPtr("Test"),
				Description: utils.StringPtr("Description"),
				Price:       utils.Float64Ptr(99.99),
				Stock:       utils.IntPtr(10),
				Page:        0,
				Size:        10,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				products := []models.Product{
					{
						ID:          uuid.New(),
						Name:        "Test Product",
						Description: utils.StringPtr("Test Description"),
						Price:       99.99,
						Stock:       10,
						CreatedAt:   time.Now(),
						UpdatedAt:   utils.TimePtr(time.Now()),
					},
				}
				repo.On("SearchProducts", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.SearchProductsRequest")).Return(products, 1, nil)
			},
			expectError:   false,
			expectedLen:   1,
			expectedTotal: 1,
		},
		{
			name: "repository error should return error",
			request: &models.SearchProductsRequest{
				Name: utils.StringPtr("Test"),
				Page: 0,
				Size: 10,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("SearchProducts", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.SearchProductsRequest")).Return(([]models.Product)(nil), 0, errors.ErrInternalServer)
			},
			expectError:   true,
			expectedLen:   0,
			expectedTotal: 0,
		},
		{
			name: "search with pagination",
			request: &models.SearchProductsRequest{
				Page: 1,
				Size: 5,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				products := []models.Product{
					{
						ID:          uuid.New(),
						Name:        "Product 6",
						Description: utils.StringPtr("Description 6"),
						Price:       99.99,
						Stock:       10,
						CreatedAt:   time.Now(),
						UpdatedAt:   utils.TimePtr(time.Now()),
					},
				}
				repo.On("SearchProducts", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.MatchedBy(func(req *models.SearchProductsRequest) bool {
					return req.Page == 1 && req.Size == 5
				})).Return(products, 20, nil) // 20 total results, showing page 1
			},
			expectError:   false,
			expectedLen:   1,
			expectedTotal: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			products, total, err := svc.SearchProducts(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, products)
				assert.Equal(t, 0, total)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, products)
				assert.Len(t, products, tt.expectedLen)
				assert.Equal(t, tt.expectedTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductsService_UpdateStock(t *testing.T) {
	setup := mocks.NewTestSetup()
	testID := uuid.New()

	tests := []struct {
		name        string
		request     *models.UpdateStockRequest
		setupMocks  func(*mocks.MockRepository)
		expectError bool
	}{
		{
			name: "successful stock update",
			request: &models.UpdateStockRequest{
				ID:    testID,
				Stock: 25,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateStock", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateStockRequest")).Return(nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.Product)
					product.ID = testID
					product.Name = "Test Product"
					product.Description = utils.StringPtr("Test Description")
					product.Price = 99.99
					product.Stock = 25
					product.CreatedAt = time.Now()
					product.UpdatedAt = utils.TimePtr(time.Now())
				})
			},
			expectError: false,
		},
		{
			name: "successful stock update to zero",
			request: &models.UpdateStockRequest{
				ID:    testID,
				Stock: 0,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateStock", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateStockRequest")).Return(nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(2).(*models.Product)
					product.ID = testID
					product.Name = "Test Product"
					product.Description = utils.StringPtr("Test Description")
					product.Price = 99.99
					product.Stock = 0
					product.CreatedAt = time.Now()
					product.UpdatedAt = utils.TimePtr(time.Now())
				})
			},
			expectError: false,
		},
		{
			name: "validation error - negative stock",
			request: &models.UpdateStockRequest{
				ID:    testID,
				Stock: -5,
			},
			setupMocks:  func(repo *mocks.MockRepository) {},
			expectError: true,
		},
		{
			name: "repository update error should return error",
			request: &models.UpdateStockRequest{
				ID:    testID,
				Stock: 25,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateStock", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateStockRequest")).Return(errors.ErrInternalServer)
			},
			expectError: true,
		},
		{
			name: "repository get error after update should return error",
			request: &models.UpdateStockRequest{
				ID:    testID,
				Stock: 25,
			},
			setupMocks: func(repo *mocks.MockRepository) {
				repo.On("UpdateStock", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateStockRequest")).Return(nil)
				repo.On("GetGeneric", setup.Ctx, mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.Product")).Return(errors.ErrNotFound)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := setup.NewMockRepository()
			tt.setupMocks(mockRepo)

			svc := service.NewService(setup.Logger, mockRepo, setup.Config)

			result, err := svc.UpdateStock(setup.Ctx, setup.Logger, tt.request)

			if tt.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testID, result.ID)
				assert.Equal(t, tt.request.Stock, result.Stock)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}