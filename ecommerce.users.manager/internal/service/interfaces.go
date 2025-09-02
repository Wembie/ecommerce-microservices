package service

import (
	"context"

	"ecommerce.users.manager/internal/models"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

type UserRepository interface {
	genericRepository
	usersRepository
}

type genericRepository interface {
	SearchGeneric(context.Context, *zap.Logger, any, int, int, func(bun.QueryBuilder) bun.QueryBuilder) (int, error)
	CreateGeneric(context.Context, *zap.Logger, any) error
	GetGeneric(context.Context, *zap.Logger, any) error
	DeleteGeneric(context.Context, *zap.Logger, any) error
	UpdateGeneric(context.Context, *zap.Logger, any) (int64, error)
}

type usersRepository interface {
	GetByUsernameOrEmail(context.Context, *zap.Logger, string) (*models.User, error)
}