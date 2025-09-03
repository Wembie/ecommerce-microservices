package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ecommerce.products.manager/internal/api/routes"
	"ecommerce.products.manager/internal/config"
	"ecommerce.products.manager/internal/repository"
	"ecommerce.products.manager/internal/service"
	_ "ecommerce.products.manager/docs"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer func() {
		_ = logger.Sync()
	}()
	logger.Info("Init products manager")

	conf := config.GetConfig()
	if conf.AppPort == "" {
		logger.Fatal("AppPort not configured")
	}

	db := config.GetConnection(conf)
	defer db.Close()

	productRepo := repository.NewRepository(db, conf)
	productService := service.NewService(logger, productRepo, conf)

	router := routes.NewRouter(logger, productService)

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