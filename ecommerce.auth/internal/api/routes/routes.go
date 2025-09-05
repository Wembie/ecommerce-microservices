package routes

import (
	"net/http"

	"ecommerce.auth/internal/api/handlers"
	middleware "ecommerce.auth/internal/api/routes/middleware"
	"ecommerce.auth/internal/clients"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, userClient clients.UserClient) *gin.Engine {
	r := gin.Default()

	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(r)

	r.Use(middleware.TimeoutMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RequestIDMiddleware())

	h := handlers.NewHandler(logger, userClient)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	auth := r.Group("/auth")
	{
		auth.POST("/token", h.Auth.GetToken)
	}

	return r
}