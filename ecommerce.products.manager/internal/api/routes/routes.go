package routes

import (
	"net/http"

	"ecommerce.products.manager/internal/api/handlers"
	middleware "ecommerce.products.manager/internal/api/routes/middleware"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, s handlers.Service) *gin.Engine {
	r := gin.Default()

	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(r)

	r.Use(middleware.TimeoutMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RequestIDMiddleware())

	h := handlers.NewHandler(logger, s)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	products := r.Group("/products")
	{
		products.GET("/search", h.Products.SearchProducts)
		products.POST("", h.Products.CreateProduct)
		products.GET("/:id", h.Products.GetProduct)
		products.PUT("/:id", h.Products.UpdateProduct)
		products.DELETE("/:id", h.Products.DeleteProduct)
		products.PUT("/:id/stock", h.Products.UpdateStock)
	}

	return r
}