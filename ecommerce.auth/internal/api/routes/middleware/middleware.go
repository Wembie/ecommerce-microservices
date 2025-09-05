package routes

import (
	"net/http"
	"strings"
	"time"

	"ecommerce.auth/internal/clients"
	"ecommerce.auth/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TimeOutResponse(c *gin.Context) {
	c.String(http.StatusRequestTimeout, "timeout")
}

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(50*time.Second),
		timeout.WithResponse(TimeOutResponse),
	)
}

func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Trace-Id", "X-Request-Id"}
	return cors.New(config)
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := config.ExtractTraceIDFromGin(c)
		c.Set("trace_id", traceID)
		c.Header("X-Trace-Id", traceID)
		c.Next()
	}
}

func BearerTokenMiddleware(userClient clients.UserClient, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("No authorization header provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn("Invalid authorization type, expected Bearer token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			logger.Warn("Empty bearer token provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token cannot be empty"})
			c.Abort()
			return
		}

		logger.Info("Validating bearer token")

		validateResp, err := userClient.ValidateUser(c.Request.Context(), token)
		if err != nil {
			logger.Error("Failed to validate token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token validation service unavailable"})
			c.Abort()
			return
		}

		if !validateResp.Valid {
			logger.Warn("Token validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}


		logger.Info("Token validated successfully", zap.String("user_id", func() string {
			if validateResp.UserID != nil {
				return validateResp.UserID.String()
			}
			return ""
		}()))

		if validateResp.UserID != nil {
			c.Set("user_id", *validateResp.UserID)
		}
		if validateResp.Username != nil {
			c.Set("username", *validateResp.Username)
		}
		if validateResp.Email != nil {
			c.Set("email", *validateResp.Email)
		}
		c.Set("token", token)
		c.Set("authenticated", true)

		c.Next()
	}
}