package handlers

import (
	"net/http"
	"strconv"

	"ecommerce.orders.manager/internal/config"
	"ecommerce.orders.manager/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrdersHandler struct {
	service Service
	logger  *zap.Logger
}

func NewOrdersHandler(logger *zap.Logger, service Service) *OrdersHandler {
	return &OrdersHandler{
		service: service,
		logger:  logger,
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with items
// @Tags Orders
// @Accept json
// @Produce json
// @Security OAuth2PasswordBearer
// @Param order body models.CreateOrderRequest true "Order data"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /orders [post]
func (h *OrdersHandler) CreateOrder(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: CreateOrder")

	var request models.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logger.Error("User ID is not a string")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	request.UserID = parsedUserID

	order, serviceErr := h.service.CreateOrder(c.Request.Context(), logger, &request)
	if serviceErr != nil {
		logger.Error("Failed to create order", zap.Error(serviceErr.Error))
		c.JSON(serviceErr.StatusCode, gin.H{"error": serviceErr.Message})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Description Get order details with items
// @Tags Orders
// @Produce json
// @Security OAuth2PasswordBearer
// @Param id path string true "Order ID"
// @Success 200 {object} models.OrderResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrdersHandler) GetOrder(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: GetOrder")

	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		logger.Error("Invalid order ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	request := models.GetOrderRequest{ID: orderID}
	order, serviceErr := h.service.GetOrder(c.Request.Context(), logger, &request)
	if serviceErr != nil {
		logger.Error("Failed to get order", zap.Error(serviceErr.Error))
		c.JSON(serviceErr.StatusCode, gin.H{"error": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrdersByUser godoc
// @Summary Get orders for authenticated user
// @Description Get paginated orders for the authenticated user
// @Tags Orders
// @Produce json
// @Security OAuth2PasswordBearer
// @Param page query int false "Page number" default(0)
// @Param size query int false "Page size" default(50)
// @Success 200 {object} docs.PaginatedResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Router /orders/user [get]
func (h *OrdersHandler) GetOrdersByUser(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: GetOrdersByUser")

	userIDValue, exists := c.Get("user_id")
	if !exists {
		logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userIDValue.(string)
	if !ok {
		logger.Error("User ID is not a string")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	//TODO improve pagination as in the products

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	request := models.GetOrdersByUserRequest{
		UserID: userID,
		Page:   page,
		Size:   size,
	}

	orders, total, serviceErr := h.service.GetOrdersByUser(c.Request.Context(), logger, &request)
	if serviceErr != nil {
		logger.Error("Failed to get orders by user", zap.Error(serviceErr.Error))
		c.JSON(serviceErr.StatusCode, gin.H{"error": serviceErr.Message})
		return
	}

	totalPages := (total + size - 1) / size

	response := gin.H{
		"data":        orders,
		"page":        page,
		"size":        size,
		"total":       total,
		"total_pages": totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an existing order
// @Tags Orders
// @Accept json
// @Produce json
// @Security OAuth2PasswordBearer
// @Param id path string true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "Status update data"
// @Success 200 {object} models.Order
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Router /orders/{id}/status [put]
func (h *OrdersHandler) UpdateOrderStatus(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: UpdateOrderStatus")

	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		logger.Error("Invalid order ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var request models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request.ID = orderID

	order, serviceErr := h.service.UpdateOrderStatus(c.Request.Context(), logger, &request)
	if serviceErr != nil {
		logger.Error("Failed to update order status", zap.Error(serviceErr.Error))
		c.JSON(serviceErr.StatusCode, gin.H{"error": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderItems godoc
// @Summary Get order items
// @Description Get all items for a specific order
// @Tags Orders
// @Produce json
// @Security OAuth2PasswordBearer
// @Param id path string true "Order ID"
// @Success 200 {array} models.OrderItem
// @Failure 400 {object} docs.ErrorResponse
// @Failure 401 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Router /orders/{id}/items [get]
func (h *OrdersHandler) GetOrderItems(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: GetOrderItems")

	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		logger.Error("Invalid order ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	items, serviceErr := h.service.GetOrderItems(c.Request.Context(), logger, orderID)
	if serviceErr != nil {
		logger.Error("Failed to get order items", zap.Error(serviceErr.Error))
		c.JSON(serviceErr.StatusCode, gin.H{"error": serviceErr.Message})
		return
	}

	c.JSON(http.StatusOK, items)
}