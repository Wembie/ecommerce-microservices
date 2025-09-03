package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ecommerce.products.manager/internal/errors"
	"ecommerce.products.manager/internal/mocks"
	"ecommerce.products.manager/internal/models"
	"ecommerce.products.manager/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupRouter(mockService *mocks.MockService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()
	return NewRouter(logger, mockService)
}

func TestNewRouter(t *testing.T) {
	mockService := &mocks.MockService{}
	logger := zap.NewNop()
	
	router := NewRouter(logger, mockService)
	
	assert.NotNil(t, router)
	
	routes := router.Routes()
	routePaths := make([]string, len(routes))
	for i, route := range routes {
		routePaths[i] = route.Path
	}
	
	expectedRoutes := []string{
		"/health",
		"/swagger/*any",
		"/products/search",
		"/products",
		"/products/:id",
		"/products/:id/stock",
	}
	
	for _, expectedRoute := range expectedRoutes {
		assert.Contains(t, routePaths, expectedRoute)
	}
}

func TestHealthEndpoint(t *testing.T) {
	mockService := &mocks.MockService{}
	router := setupRouter(mockService)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    models.CreateProductRequest
		mockSetup      func(*mocks.MockService)
		expectedStatus int
	}{
		{
			name: "successful product creation",
			requestBody: models.CreateProductRequest{
				Name:        "Test Product",
				Description: utils.StringPtr("Test Description"),
				Price:       29.99,
				Stock:       100,
			},
			mockSetup: func(ms *mocks.MockService) {
				product := &models.Product{
					ID:          uuid.New(),
					Name:        "Test Product",
					Description: utils.StringPtr("Test Description"),
					Price:       29.99,
					Stock:       100,
					CreatedAt:   time.Now(),
				}
				ms.On("CreateProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.CreateProductRequest")).Return(product, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid request body",
			requestBody: models.CreateProductRequest{
				Name:  "",
				Price: -1,
				Stock: -1,
			},
			mockSetup: func(ms *mocks.MockService) {
				errorAPI := errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "validation error")
				ms.On("CreateProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.CreateProductRequest")).Return((*models.Product)(nil), errorAPI)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockService{}
			tt.mockSetup(mockService)
			router := setupRouter(mockService)
			
			jsonBody, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetProduct(t *testing.T) {
	productID := uuid.New()
	
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*mocks.MockService)
		expectedStatus int
	}{
		{
			name:      "successful product retrieval",
			productID: productID.String(),
			mockSetup: func(ms *mocks.MockService) {
				product := &models.Product{
					ID:          productID,
					Name:        "Test Product",
					Description: utils.StringPtr("Test Description"),
					Price:       29.99,
					Stock:       100,
					CreatedAt:   time.Now(),
				}
				ms.On("GetProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.GetProductRequest")).Return(product, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "product not found",
			productID: productID.String(),
			mockSetup: func(ms *mocks.MockService) {
				errorAPI := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "product not found")
				ms.On("GetProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.GetProductRequest")).Return((*models.Product)(nil), errorAPI)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid UUID",
			productID:      "invalid-uuid",
			mockSetup:      func(ms *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockService{}
			tt.mockSetup(mockService)
			router := setupRouter(mockService)
			
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/products/%s", tt.productID), nil)
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	productID := uuid.New()
	
	tests := []struct {
		name           string
		productID      string
		requestBody    models.UpdateProductRequest
		mockSetup      func(*mocks.MockService)
		expectedStatus int
	}{
		{
			name:      "successful product update",
			productID: productID.String(),
			requestBody: models.UpdateProductRequest{
				Name:  utils.StringPtr("Updated Product"),
				Price: utils.Float64Ptr(39.99),
			},
			mockSetup: func(ms *mocks.MockService) {
				updatedProduct := &models.Product{
					ID:        productID,
					Name:      "Updated Product",
					Price:     39.99,
					Stock:     100,
					CreatedAt: time.Now(),
					UpdatedAt: utils.TimePtr(time.Now()),
				}
				ms.On("UpdateProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateProductRequest")).Return(updatedProduct, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "product not found",
			productID: productID.String(),
			requestBody: models.UpdateProductRequest{
				Name: utils.StringPtr("Updated Product"),
			},
			mockSetup: func(ms *mocks.MockService) {
				errorAPI := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "product not found")
				ms.On("UpdateProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateProductRequest")).Return((*models.Product)(nil), errorAPI)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid UUID",
			productID:      "invalid-uuid",
			requestBody:    models.UpdateProductRequest{},
			mockSetup:      func(ms *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockService{}
			tt.mockSetup(mockService)
			router := setupRouter(mockService)
			
			jsonBody, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%s", tt.productID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	productID := uuid.New()
	
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*mocks.MockService)
		expectedStatus int
	}{
		{
			name:      "successful product deletion",
			productID: productID.String(),
			mockSetup: func(ms *mocks.MockService) {
				response := &models.DeleteProductResponse{Success: true}
				ms.On("DeleteProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.DeleteProductRequest")).Return(response, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "product not found",
			productID: productID.String(),
			mockSetup: func(ms *mocks.MockService) {
				errorAPI := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "product not found")
				ms.On("DeleteProduct", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.DeleteProductRequest")).Return((*models.DeleteProductResponse)(nil), errorAPI)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid UUID",
			productID:      "invalid-uuid",
			mockSetup:      func(ms *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockService{}
			tt.mockSetup(mockService)
			router := setupRouter(mockService)
			
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/products/%s", tt.productID), nil)
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateStock(t *testing.T) {
	productID := uuid.New()
	
	tests := []struct {
		name           string
		productID      string
		requestBody    models.UpdateStockRequest
		mockSetup      func(*mocks.MockService)
		expectedStatus int
	}{
		{
			name:      "successful stock update",
			productID: productID.String(),
			requestBody: models.UpdateStockRequest{
				Stock: 50,
			},
			mockSetup: func(ms *mocks.MockService) {
				updatedProduct := &models.Product{
					ID:        productID,
					Name:      "Test Product",
					Price:     29.99,
					Stock:     50,
					CreatedAt: time.Now(),
					UpdatedAt: utils.TimePtr(time.Now()),
				}
				ms.On("UpdateStock", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateStockRequest")).Return(updatedProduct, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "product not found",
			productID: productID.String(),
			requestBody: models.UpdateStockRequest{
				Stock: 50,
			},
			mockSetup: func(ms *mocks.MockService) {
				errorAPI := errors.NewErrorAPI(http.StatusNotFound, errors.ErrNotFound, "product not found")
				ms.On("UpdateStock", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.UpdateStockRequest")).Return((*models.Product)(nil), errorAPI)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid UUID",
			productID:      "invalid-uuid",
			requestBody:    models.UpdateStockRequest{Stock: 50},
			mockSetup:      func(ms *mocks.MockService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockService{}
			tt.mockSetup(mockService)
			router := setupRouter(mockService)
			
			jsonBody, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/products/%s/stock", tt.productID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSearchProducts(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*mocks.MockService)
		expectedStatus int
	}{
		{
			name:        "successful search with name filter",
			queryParams: "?name=test&page=1&size=10",
			mockSetup: func(ms *mocks.MockService) {
				products := []models.Product{
					{
						ID:          uuid.New(),
						Name:        "Test Product 1",
						Description: utils.StringPtr("Description 1"),
						Price:       29.99,
						Stock:       100,
						CreatedAt:   time.Now(),
					},
				}
				ms.On("SearchProducts", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.SearchProductsRequest")).Return(products, 1, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "successful search with multiple filters",
			queryParams: "?name=test&price=29.99&stock=100&page=1&size=10",
			mockSetup: func(ms *mocks.MockService) {
				products := []models.Product{}
				ms.On("SearchProducts", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.SearchProductsRequest")).Return(products, 0, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "search with no results",
			queryParams: "?name=nonexistent&page=1&size=10",
			mockSetup: func(ms *mocks.MockService) {
				products := []models.Product{}
				ms.On("SearchProducts", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("*zap.Logger"), mock.AnythingOfType("*models.SearchProductsRequest")).Return(products, 0, (*errors.ErrorAPI)(nil))
			},
			expectedStatus: http.StatusOK,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mocks.MockService{}
			tt.mockSetup(mockService)
			router := setupRouter(mockService)
			
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/products/search"+tt.queryParams, nil)
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestSwaggerEndpoint(t *testing.T) {
	mockService := &mocks.MockService{}
	router := setupRouter(mockService)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/swagger/index.html", nil)
	router.ServeHTTP(w, req)
	
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)
}

func TestMetricsEndpoint(t *testing.T) {
	mockService := &mocks.MockService{}
	router := setupRouter(mockService)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
}