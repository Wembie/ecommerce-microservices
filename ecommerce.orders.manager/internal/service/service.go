package service

import (
	"ecommerce.orders.manager/internal/api/handlers"
	"ecommerce.orders.manager/internal/clients"
	"ecommerce.orders.manager/internal/config"

	"go.uber.org/zap"
)

type OrdersService struct {
	orderRepo      OrderRepository
	userClient     clients.UserClient
	productClient  clients.ProductClient
	log           *zap.Logger
	conf          *config.Config
}

var _ handlers.Service = (*OrdersService)(nil)

func NewService(
	log *zap.Logger,
	repo OrderRepository,
	userClient clients.UserClient,
	productClient clients.ProductClient,
	conf *config.Config,
) *OrdersService {
	return &OrdersService{
		orderRepo:     repo,
		userClient:    userClient,
		productClient: productClient,
		log:          log,
		conf:         conf,
	}
}