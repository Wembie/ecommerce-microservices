package transport

import (
	"ecommerce.users.manager/internal/config"

	pb "github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo"
	"go.uber.org/zap"
)

type PingHandler struct {
	pb.UserServiceServer
	log     *zap.Logger
	pingSvc Service
	cfg     *config.Config
}

func NewHandler(log *zap.Logger, pingSvc Service, cfg *config.Config) *PingHandler {
	return &PingHandler{
		log:     log,
		pingSvc: pingSvc,
		cfg:     cfg,
	}
}
