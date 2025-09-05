package handlers

import (
	"net/http"

	"ecommerce.orders.manager/internal/clients"
	"ecommerce.orders.manager/internal/config"
	"ecommerce.orders.manager/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userClient clients.UserClient
	logger     *zap.Logger
}

func NewAuthHandler(logger *zap.Logger, userClient clients.UserClient) *AuthHandler {
	return &AuthHandler{
		userClient: userClient,
		logger:     logger,
	}
}

// GetToken godoc
// @Summary Get access token
// @Description OAuth2 password flow - exchange username/password for access token
// @Tags Authentication
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/token [post]
func (h *AuthHandler) GetToken(c *gin.Context) {
	logger := config.CreateLoggerWithTraceID(h.logger, c)
	logger.Info("Handler: GetToken")

	var request models.TokenRequest
	if err := c.ShouldBind(&request); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	authResp, err := h.userClient.AuthenticateUser(c.Request.Context(), &clients.AuthRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		logger.Error("Failed to authenticate user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service unavailable"})
		return
	}

	if !authResp.Success {
		logger.Info("Authentication failed", zap.String("username", request.Username))
		errorMsg := "Invalid username or password"
		if authResp.ErrorMessage != nil {
			errorMsg = *authResp.ErrorMessage
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
		return
	}

	if authResp.Token == nil {
		logger.Error("No token received from authentication service")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication service error"})
		return
	}

	logger.Info("Token generated successfully", zap.String("username", request.Username))
	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken: *authResp.Token,
		TokenType:   "bearer",
	})
}