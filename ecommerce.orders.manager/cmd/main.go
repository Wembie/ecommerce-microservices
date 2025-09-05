package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ecommerce.orders.manager/internal/api/routes"
	"ecommerce.orders.manager/internal/clients"
	"ecommerce.orders.manager/internal/config"
	"ecommerce.orders.manager/internal/repository"
	"ecommerce.orders.manager/internal/service"
	_ "ecommerce.orders.manager/docs"

	"go.uber.org/zap"
)

// @title Orders Manager API
// @version 1.0
// @description Orders Manager service for e-commerce microservices
// @host localhost:8082
// @BasePath /
// @securityDefinitions.oauth2.password OAuth2PasswordBearer
// @tokenUrl /auth/token

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer func() {
		_ = logger.Sync()
	}()
	logger.Info("Init orders manager")

	conf := config.GetConfig()
	if conf.AppPort == "" {
		logger.Fatal("AppPort not configured")
	}

	db := config.GetConnection(conf)
	defer db.Close()

	userClient, err := clients.NewUserClient(conf.UserManagerHost, logger)
	if err != nil {
		logger.Fatal("Failed to create user client", zap.Error(err))
	}
	defer userClient.Close()

	productClient := clients.NewProductClient(conf.ProductManagerHost, logger)

	orderRepo := repository.NewRepository(db, conf)
	orderService := service.NewService(logger, orderRepo, userClient, productClient, conf)

	router := routes.NewRouter(logger, orderService, userClient)

	go func() {
		logger.Info("Starting HTTP server", zap.String("port", conf.AppPort))
		if err := router.Run(":" + conf.AppPort); err != nil {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	logger.Info("Server stopped cleanly")
}