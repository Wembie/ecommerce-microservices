package routes

import (
	"net/http"

	"ecommerce.orders.manager/internal/api/handlers"
	middleware "ecommerce.orders.manager/internal/api/routes/middleware"
	"ecommerce.orders.manager/internal/clients"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, s handlers.Service, userClient clients.UserClient) *gin.Engine {
	r := gin.Default()

	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(r)

	r.Use(middleware.TimeoutMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RequestIDMiddleware())

	h := handlers.NewHandler(logger, s, userClient)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	{
		auth.POST("/token", h.Auth.GetToken)
	}

	protected := r.Group("/")
	protected.Use(middleware.BearerTokenMiddleware(userClient, logger))
	{
		orders := protected.Group("/orders")
		{
			orders.POST("", h.Orders.CreateOrder)
			orders.GET("/:id", h.Orders.GetOrder)
			orders.GET("/user", h.Orders.GetOrdersByUser)
			orders.PUT("/:id/status", h.Orders.UpdateOrderStatus)
			orders.GET("/:id/items", h.Orders.GetOrderItems)
		}
	}

	return r
}
