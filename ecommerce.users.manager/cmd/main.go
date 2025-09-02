package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"ecommerce.users.manager/internal/config"
	"ecommerce.users.manager/internal/repository"
	"ecommerce.users.manager/internal/service"
	"ecommerce.users.manager/internal/transport"

	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer func() {
		_ = logger.Sync()
	}()
	logger.Info("Init users manager")

	conf := config.GetConfig()
	if conf.AppPort == "" {
		logger.Fatal("AppPort not configured")
	}

	db := config.GetConnection(conf)
	defer db.Close()

	userRepo := repository.NewRepository(db, conf)
	userService := service.NewService(logger, userRepo, conf)
	userHandler := transport.NewHandler(logger, userService, conf)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", conf.AppPort))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}
	logger.Info("gRPC server listening", zap.String("address", lis.Addr().String()))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	logger.Info("Server stopped cleanly")
}