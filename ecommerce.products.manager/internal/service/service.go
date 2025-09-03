package service

import (
	"ecommerce.products.manager/internal/config"
	"ecommerce.products.manager/internal/api/handlers"

	"go.uber.org/zap"
)

type ProductsService struct {
	productRepo ProductRepository
	log        *zap.Logger
	conf       *config.Config
}

var _ handlers.Service = (*ProductsService)(nil)

func NewService(log *zap.Logger, repo ProductRepository, conf *config.Config) *ProductsService {
	return &ProductsService{
		productRepo: repo,
		log:        log,
		conf:       conf,
	}
}