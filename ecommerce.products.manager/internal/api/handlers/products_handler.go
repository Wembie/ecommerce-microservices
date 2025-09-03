package handlers

import (
	"net/http"
	"strconv"

	"ecommerce.products.manager/internal/config"
	"ecommerce.products.manager/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductsHandler struct {
	service Service
	logger  *zap.Logger
}

func NewProductsHandler(logger *zap.Logger, service Service) *ProductsHandler {
	return &ProductsHandler{
		service: service,
		logger:  logger,
	}
}


// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Product data"
// @Success 201 {object} models.Product
// @Failure 400 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /products [post]
func (h *ProductsHandler) CreateProduct(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: CreateProduct")

	var request models.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), logger, &request)
	if err != nil {
		logger.Error("Failed to create product", zap.Error(err.Error))
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductsHandler) GetProduct(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: GetProduct")

	idStr := c.Param("id")
	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		logger.Error("Invalid UUID", zap.Error(parseErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	request := &models.GetProductRequest{ID: id}
	product, err := h.service.GetProduct(c.Request.Context(), logger, request)
	if err != nil {
		logger.Error("Failed to get product", zap.Error(err.Error))
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.UpdateProductRequest true "Product data"
// @Success 200 {object} models.Product
// @Failure 400 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /products/{id} [put]
func (h *ProductsHandler) UpdateProduct(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: UpdateProduct")

	idStr := c.Param("id")
	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		logger.Error("Invalid UUID", zap.Error(parseErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var request models.UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.ID = id
	product, err := h.service.UpdateProduct(c.Request.Context(), logger, &request)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err.Error))
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.DeleteProductResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductsHandler) DeleteProduct(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: DeleteProduct")

	idStr := c.Param("id")
	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		logger.Error("Invalid UUID", zap.Error(parseErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	request := &models.DeleteProductRequest{ID: id}
	response, err := h.service.DeleteProduct(c.Request.Context(), logger, request)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err.Error))
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SearchProducts godoc
// @Summary Search and list products with advanced filters
// @Description Search products using multiple criteria including name, description, price, and stock filters with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param name query string false "Filter by product name (partial match)"
// @Param description query string false "Filter by product description (partial match)"
// @Param price query number false "Filter by exact price"
// @Param stock query integer false "Filter by minimum stock quantity"
// @Param page query integer false "Page number (0-based)" default(0)
// @Param size query integer false "Number of items per page" default(50)
// @Success 200 {object} docs.PaginatedProductResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /products/search [get]
func (h *ProductsHandler) SearchProducts(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: SearchProducts")

	var request models.SearchProductsRequest

	if name := c.Query("name"); name != "" {
		request.Name = &name
	}

	if description := c.Query("description"); description != "" {
		request.Description = &description
	}

	if priceStr := c.Query("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err != nil {
			logger.Warn("Invalid price parameter", zap.String("price", priceStr), zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format: must be a number"})
			return
		} else if price < 0 {
			logger.Warn("Invalid price filter", zap.Float64("price", price))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Price filter must be non-negative"})
			return
		} else {
			request.Price = &price
		}
	}

	if stockStr := c.Query("stock"); stockStr != "" {
		if stock, err := strconv.Atoi(stockStr); err != nil {
			logger.Warn("Invalid stock parameter", zap.String("stock", stockStr), zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock format: must be an integer"})
			return
		} else if stock < 0 {
			logger.Warn("Invalid stock filter", zap.Int("stock", stock))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock filter must be non-negative"})
			return
		} else {
			request.Stock = &stock
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "50"))

	if page < 0 {
		page = 0
	}
	if size <= 0 || size > 100 {
		size = 50
	}

	request.Page = page
	request.Size = size

	products, total, err := h.service.SearchProducts(c.Request.Context(), logger, &request)
	if err != nil {
		logger.Error("Failed to search products", zap.Error(err.Error))
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	response := models.NewPaginatedResponse(products, request.Page, request.Size, total)
	c.JSON(http.StatusOK, response)
}

// UpdateStock godoc
// @Summary Update product stock
// @Description Update the stock of a product
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param stock body models.UpdateStockRequest true "Stock data"
// @Success 200 {object} models.Product
// @Failure 400 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /products/{id}/stock [put]
func (h *ProductsHandler) UpdateStock(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: UpdateStock")

	idStr := c.Param("id")
	id, parseErr := uuid.Parse(idStr)
	if parseErr != nil {
		logger.Error("Invalid UUID", zap.Error(parseErr))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var request models.UpdateStockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.ID = id
	product, err := h.service.UpdateStock(c.Request.Context(), logger, &request)
	if err != nil {
		logger.Error("Failed to update stock", zap.Error(err.Error))
		c.JSON(err.StatusCode, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, product)
}