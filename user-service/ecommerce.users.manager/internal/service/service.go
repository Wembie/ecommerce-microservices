package service

import (
	"ecommerce.users.manager/internal/config"
	"ecommerce.users.manager/internal/transport"

	"go.uber.org/zap"
)

type UsersService struct {
	userRepo	UserRepository
	log         *zap.Logger
	conf        *config.Config
	jwtKey      []byte
}

var _ transport.Service = (*UsersService)(nil)

func NewService(log *zap.Logger, repo UserRepository, conf *config.Config) *UsersService {
	return &UsersService{
		userRepo:	repo,
		log:        log,
		conf:       conf,
		jwtKey:     []byte(conf.JWTSecret),
	}
}
