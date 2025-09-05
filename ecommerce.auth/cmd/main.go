package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ecommerce.auth/internal/api/routes"
	"ecommerce.auth/internal/clients"
	"ecommerce.auth/internal/config"
	_ "ecommerce.auth/docs"

	"go.uber.org/zap"
)

// @title           E-commerce Auth API
// @version         1.0
// @description     Authentication service for e-commerce microservices
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer func() {
		_ = logger.Sync()
	}()
	logger.Info("Init auth manager")

	conf := config.GetConfig()
	if conf.AppPort == "" {
		logger.Fatal("AppPort not configured")
	}

	userClient, err := clients.NewUserClient(conf.UserManagerHost, logger)
	if err != nil {
		logger.Fatal("Failed to create user client", zap.Error(err))
	}
	defer userClient.Close()

	router := routes.NewRouter(logger, userClient)

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
