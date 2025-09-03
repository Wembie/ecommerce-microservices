package routes

import (
	"net/http"
	"time"

	"ecommerce.products.manager/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
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