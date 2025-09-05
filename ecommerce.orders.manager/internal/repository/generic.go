package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	internalErrors "ecommerce.orders.manager/internal/errors"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

func (r *orderRepository) CreateGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	logger.Info("Repository: CreateGeneric")

	if reflect.ValueOf(model).Kind().String() != "ptr" && reflect.ValueOf(model).Kind().String() != "struct" {
		return internalErrors.ErrInvalidInput
	}

	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to create record", zap.Error(err))
		return err
	}

	return nil
}

func (r *orderRepository) GetGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	logger.Info("Repository: GetGeneric")

	if reflect.ValueOf(model).Kind().String() != "ptr" && reflect.ValueOf(model).Kind().String() != "struct" {
		return internalErrors.ErrInvalidInput
	}

	err := r.db.NewSelect().
		Model(model).
		WherePK().
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("Record not found")
			return internalErrors.ErrNotFound
		}
		logger.Error("Failed to get record", zap.Error(err))
		return err
	}

	return nil
}

func (r *orderRepository) UpdateGeneric(ctx context.Context, logger *zap.Logger, model any) (int64, error) {
	logger.Info("Repository: UpdateGeneric")

	if reflect.ValueOf(model).Kind().String() != "ptr" && reflect.ValueOf(model).Kind().String() != "struct" {
		return 0, internalErrors.ErrInvalidInput
	}

	result, err := r.db.NewUpdate().
		Model(model).
		WherePK().
		OmitZero().
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to update record", zap.Error(err))
		return 0, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("No rows updated, record may not exist")
		return 0, internalErrors.ErrNotFound
	}

	return rowsAffected, nil
}

func (r *orderRepository) DeleteGeneric(ctx context.Context, logger *zap.Logger, model any) error {
	logger.Info("Repository: DeleteGeneric")

	if reflect.ValueOf(model).Kind().String() != "ptr" && reflect.ValueOf(model).Kind().String() != "struct" {
		return internalErrors.ErrInvalidInput
	}

	result, err := r.db.NewDelete().
		Model(model).
		WherePK().
		Exec(ctx)

	if err != nil {
		logger.Error("Failed to delete record", zap.Error(err))
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		logger.Warn("No rows deleted, record may not exist")
		return internalErrors.ErrNotFound
	}

	return nil
}

func (r *orderRepository) SearchGeneric(
	ctx context.Context,
	logger *zap.Logger,
	model any,
	page int,
	size int,
	queryModifier func(*bun.SelectQuery) *bun.SelectQuery,
) (int, error) {
	logger.Info("Repository: SearchGeneric", zap.Int("page", page), zap.Int("size", size))

	if reflect.ValueOf(model).Kind().String() != "ptr" && reflect.ValueOf(model).Kind().String() != "struct" {
		return 0, internalErrors.ErrInvalidInput
	}
	offset := (page - 1) * size

	query := r.db.NewSelect().Model(model)

	if queryModifier != nil {
		query = queryModifier(query)
	}

	count, err := query.
		Limit(size).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("No records found")
			return 0, internalErrors.ErrNotFound
		}
		logger.Error("Failed to search records", zap.Error(err))
		return 0, err
	}

	return count, nil
}