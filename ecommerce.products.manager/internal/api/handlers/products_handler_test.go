package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ecommerce.products.manager/internal/errors"
	"ecommerce.products.manager/internal/api/handlers"
	"ecommerce.products.manager/internal/mocks"
	"ecommerce.products.manager/internal/models"
	"ecommerce.products.manager/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductsHandler_CreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup := mocks.NewTestSetup()

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mocks.MockService)
		expectedStatus int
		expectError    bool
	}{
		{
			name: "successful create",
			requestBody: models.CreateProductRequest{
				Name:        "Test Product",
				Description: utils.StringPtr("Test Description"),
				Price:       99.99,
				Stock:       10,
			},
			setupMocks: func(svc *mocks.MockService) {
				product := &models.Product{
					ID:          uuid.New(),
					Name:        "Test Product",
					Description: utils.StringPtr("Test Description"),
					Price:       99.99,
					Stock:       10,
					CreatedAt:   time.Now(),
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
				svc.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return(product, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "service error",
			requestBody: models.CreateProductRequest{
				Name:        "Test Product",
				Description: utils.StringPtr("Test Description"),
				Price:       99.99,
				Stock:       10,
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusInternalServerError, errors.ErrInternalServer, "Service error")
				svc.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name: "validation error - duplicate resource",
			requestBody: models.CreateProductRequest{
				Name:        "Existing Product",
				Description: utils.StringPtr("Test Description"),
				Price:       99.99,
				Stock:       10,
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusConflict, errors.ErrDuplicateResource, "Product already exists")
				svc.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusConflict,
			expectError:    true,
		},
		{
			name: "validation error - invalid input",
			requestBody: models.CreateProductRequest{
				Name:        "",
				Description: utils.StringPtr("Test Description"),
				Price:       -1,
				Stock:       -5,
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Invalid product data")
				svc.On("CreateProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			handler := handlers.NewProductsHandler(setup.Logger, mockSvc)
			tt.setupMocks(mockSvc)

			router := gin.New()
			router.POST("/products", handler.CreateProduct)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				var response models.Product
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, "Test Product", response.Name)
			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestProductsHandler_GetProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup := mocks.NewTestSetup()
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		setupMocks     func(*mocks.MockService)
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful get",
			productID: validID.String(),
			setupMocks: func(svc *mocks.MockService) {
				product := &models.Product{
					ID:          validID,
					Name:        "Test Product",
					Description: utils.StringPtr("Test Description"),
					Price:       99.99,
					Stock:       10,
					CreatedAt:   time.Now(),
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
				svc.On("GetProduct", mock.Anything, mock.Anything, mock.Anything).Return(product, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid uuid",
			productID:      "invalid-uuid",
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "Product not found")
				svc.On("GetProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:      "service error",
			productID: validID.String(),
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusInternalServerError, errors.ErrInternalServer, "Database error")
				svc.On("GetProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			handler := handlers.NewProductsHandler(setup.Logger, mockSvc)
			tt.setupMocks(mockSvc)

			router := gin.New()
			router.GET("/products/:id", handler.GetProduct)

			req, _ := http.NewRequest("GET", "/products/"+tt.productID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				var response models.Product
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.productID, response.ID.String())
			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestProductsHandler_UpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup := mocks.NewTestSetup()
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		requestBody    interface{}
		setupMocks     func(*mocks.MockService)
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful update",
			productID: validID.String(),
			requestBody: models.UpdateProductRequest{
				Name:        utils.StringPtr("Updated Product"),
				Description: utils.StringPtr("Updated Description"),
				Price:       utils.Float64Ptr(149.99),
				Stock:       utils.IntPtr(5),
			},
			setupMocks: func(svc *mocks.MockService) {
				product := &models.Product{
					ID:          validID,
					Name:        "Updated Product",
					Description: utils.StringPtr("Updated Description"),
					Price:       149.99,
					Stock:       5,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
				svc.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(product, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid uuid",
			productID:      "invalid-uuid",
			requestBody:    models.UpdateProductRequest{},
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "invalid request body",
			productID:      validID.String(),
			requestBody:    "invalid json",
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			requestBody: models.UpdateProductRequest{
				Name: utils.StringPtr("Updated Product"),
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "Product not found")
				svc.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:      "invalid input error",
			productID: validID.String(),
			requestBody: models.UpdateProductRequest{
				Price: utils.Float64Ptr(-10),
				Stock: utils.IntPtr(-5),
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Price and stock must be positive")
				svc.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			handler := handlers.NewProductsHandler(setup.Logger, mockSvc)
			tt.setupMocks(mockSvc)

			router := gin.New()
			router.PUT("/products/:id", handler.UpdateProduct)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest("PUT", "/products/"+tt.productID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				var response models.Product
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.productID, response.ID.String())
			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestProductsHandler_DeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup := mocks.NewTestSetup()
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		setupMocks     func(*mocks.MockService)
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful delete",
			productID: validID.String(),
			setupMocks: func(svc *mocks.MockService) {
				response := &models.DeleteProductResponse{
					Success: true,
				}
				svc.On("DeleteProduct", mock.Anything, mock.Anything, mock.Anything).Return(response, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid uuid",
			productID:      "invalid-uuid",
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "Product not found")
				svc.On("DeleteProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.DeleteProductResponse)(nil), apiErr)
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:      "service error",
			productID: validID.String(),
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusInternalServerError, errors.ErrInternalServer, "Failed to delete product")
				svc.On("DeleteProduct", mock.Anything, mock.Anything, mock.Anything).Return((*models.DeleteProductResponse)(nil), apiErr)
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:      "delete failed - returns false",
			productID: validID.String(),
			setupMocks: func(svc *mocks.MockService) {
				response := &models.DeleteProductResponse{
					Success: false,
				}
				svc.On("DeleteProduct", mock.Anything, mock.Anything, mock.Anything).Return(response, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			handler := handlers.NewProductsHandler(setup.Logger, mockSvc)
			tt.setupMocks(mockSvc)

			router := gin.New()
			router.DELETE("/products/:id", handler.DeleteProduct)

			req, _ := http.NewRequest("DELETE", "/products/"+tt.productID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				var response models.DeleteProductResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestProductsHandler_SearchProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup := mocks.NewTestSetup()

	tests := []struct {
		name           string
		queryParams    map[string]string
		setupMocks     func(*mocks.MockService)
		expectedStatus int
		expectError    bool
	}{
		{
			name:        "successful search without filters",
			queryParams: map[string]string{},
			setupMocks: func(svc *mocks.MockService) {
				products := []models.Product{
					{
						ID:          uuid.New(),
						Name:        "Product 1",
						Description: utils.StringPtr("Description 1"),
						Price:       99.99,
						Stock:       10,
						CreatedAt:   time.Now(),
						UpdatedAt:   utils.TimePtr(time.Now()),
					},
					{
						ID:          uuid.New(),
						Name:        "Product 2",
						Description: utils.StringPtr("Description 2"),
						Price:       149.99,
						Stock:       5,
						CreatedAt:   time.Now(),
						UpdatedAt:   utils.TimePtr(time.Now()),
					},
				}
				svc.On("SearchProducts", mock.Anything, mock.Anything, mock.Anything).Return(products, 2, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "successful search with filters",
			queryParams: map[string]string{
				"name":        "Test",
				"description": "Description",
				"price":       "99.99",
				"stock":       "5",
				"page":        "0",
				"size":        "10",
			},
			setupMocks: func(svc *mocks.MockService) {
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
				svc.On("SearchProducts", mock.Anything, mock.Anything, mock.Anything).Return(products, 1, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "invalid price parameter",
			queryParams: map[string]string{
				"price": "invalid-price",
			},
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "negative price parameter",
			queryParams: map[string]string{
				"price": "-10.50",
			},
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid stock parameter",
			queryParams: map[string]string{
				"stock": "invalid-stock",
			},
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "negative stock parameter",
			queryParams: map[string]string{
				"stock": "-5",
			},
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "pagination limits - large size",
			queryParams: map[string]string{
				"page": "0",
				"size": "200", // Should be capped at 100
			},
			setupMocks: func(svc *mocks.MockService) {
				products := []models.Product{}
				svc.On("SearchProducts", mock.Anything, mock.Anything, mock.MatchedBy(func(req *models.SearchProductsRequest) bool {
					return req.Size == 50 // Should be set to default 50
				})).Return(products, 0, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "service error",
			queryParams: map[string]string{
				"name": "Test",
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusInternalServerError, errors.ErrInternalServer, "Database error")
				svc.On("SearchProducts", mock.Anything, mock.Anything, mock.Anything).Return(([]models.Product)(nil), 0, apiErr)
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name: "empty results",
			queryParams: map[string]string{
				"name": "NonExistentProduct",
			},
			setupMocks: func(svc *mocks.MockService) {
				products := []models.Product{}
				svc.On("SearchProducts", mock.Anything, mock.Anything, mock.Anything).Return(products, 0, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			handler := handlers.NewProductsHandler(setup.Logger, mockSvc)
			tt.setupMocks(mockSvc)

			router := gin.New()
			router.GET("/products/search", handler.SearchProducts)

			req, _ := http.NewRequest("GET", "/products/search", nil)
			
			// Add query parameters
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				var response models.PaginatedResponse[models.Product]
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Items)
				assert.GreaterOrEqual(t, response.Page, 0)
				assert.Greater(t, response.Size, 0)
				assert.GreaterOrEqual(t, response.Total, 0)
			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestProductsHandler_UpdateStock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setup := mocks.NewTestSetup()
	validID := uuid.New()

	tests := []struct {
		name           string
		productID      string
		requestBody    interface{}
		setupMocks     func(*mocks.MockService)
		expectedStatus int
		expectError    bool
	}{
		{
			name:      "successful stock update",
			productID: validID.String(),
			requestBody: models.UpdateStockRequest{
				Stock: 25,
			},
			setupMocks: func(svc *mocks.MockService) {
				product := &models.Product{
					ID:          validID,
					Name:        "Test Product",
					Description: utils.StringPtr("Test Description"),
					Price:       99.99,
					Stock:       25,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
				svc.On("UpdateStock", mock.Anything, mock.Anything, mock.Anything).Return(product, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid uuid",
			productID:      "invalid-uuid",
			requestBody:    models.UpdateStockRequest{Stock: 25},
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "invalid request body",
			productID:      validID.String(),
			requestBody:    "invalid json",
			setupMocks:     func(svc *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "product not found",
			productID: validID.String(),
			requestBody: models.UpdateStockRequest{
				Stock: 25,
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "Product not found")
				svc.On("UpdateStock", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:      "negative stock validation",
			productID: validID.String(),
			requestBody: models.UpdateStockRequest{
				Stock: -5,
			},
			setupMocks: func(svc *mocks.MockService) {
				apiErr := errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Stock cannot be negative")
				svc.On("UpdateStock", mock.Anything, mock.Anything, mock.Anything).Return((*models.Product)(nil), apiErr)
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:      "zero stock update",
			productID: validID.String(),
			requestBody: models.UpdateStockRequest{
				Stock: 0,
			},
			setupMocks: func(svc *mocks.MockService) {
				product := &models.Product{
					ID:          validID,
					Name:        "Test Product",
					Description: utils.StringPtr("Test Description"),
					Price:       99.99,
					Stock:       0,
					CreatedAt:   time.Now().Add(-24 * time.Hour),
					UpdatedAt:   utils.TimePtr(time.Now()),
				}
				svc.On("UpdateStock", mock.Anything, mock.Anything, mock.Anything).Return(product, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := setup.NewMockService()
			handler := handlers.NewProductsHandler(setup.Logger, mockSvc)
			tt.setupMocks(mockSvc)

			router := gin.New()
			router.PUT("/products/:id/stock", handler.UpdateStock)

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, _ := http.NewRequest("PUT", "/products/"+tt.productID+"/stock", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if !tt.expectError {
				var response models.Product
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.productID, response.ID.String())
			} else {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Contains(t, errorResponse, "error")
			}

			mockSvc.AssertExpectations(t)
		})
	}
}