package repository

import (
	"context"
	"reflect"

	"github.com/uptrace/bun"
	"go.uber.org/zap"

	"ecommerce.users.manager/internal/errors"
)

func (r *userRepository) SearchGeneric(
	ctx context.Context,
	logger *zap.Logger,
	t any,
	limit int,
	offset int,
	searchBy func(q bun.QueryBuilder) bun.QueryBuilder,
) (int, error) {
	if reflect.ValueOf(t).Kind().String() != "ptr" && reflect.ValueOf(t).Kind().String() != "struct" {
		return 0, errors.New(errors.ErrCodeInvalidArgument, "the argument must be a struct or pointer")
	}

	total, err := r.db.NewSelect().
		Model(t).
		Limit(limit).
		Offset(offset).
		ApplyQueryBuilder(searchBy).
		ScanAndCount(ctx)
	if err != nil {
		return 0, errors.New(errors.ErrCodeInternal, err.Error())
	}

	return total, nil
}

func (r *userRepository) CreateGeneric(ctx context.Context, logger *zap.Logger, t any) error {
	if reflect.ValueOf(t).Kind().String() != "ptr" && reflect.ValueOf(t).Kind().String() != "struct" {
		return errors.New(errors.ErrCodeInvalidArgument, "the argument must be a struct or pointer")
	}

	_, err := r.db.NewInsert().
		Model(t).
		Exec(ctx)
	if err != nil {
		return errors.New(errors.ErrCodeInternal, err.Error())
	}

	return nil
}

func (r *userRepository) GetGeneric(ctx context.Context, logger *zap.Logger, t any) error {
	if reflect.ValueOf(t).Kind().String() != "ptr" && reflect.ValueOf(t).Kind().String() != "struct" {
		return errors.New(errors.ErrCodeInvalidArgument, "the argument must be a struct or pointer")
	}

	err := r.db.NewSelect().
		Model(t).
		WherePK().
		Scan(ctx)
	if err != nil {
		return errors.New(errors.ErrCodeNotFound, err.Error())
	}

	return nil
}

func (r *userRepository) DeleteGeneric(ctx context.Context, logger *zap.Logger, t any) error {
	if reflect.ValueOf(t).Kind().String() != "ptr" && reflect.ValueOf(t).Kind().String() != "struct" {
		return errors.New(errors.ErrCodeInvalidArgument, "the argument must be a struct or pointer")
	}

	resp, err := r.db.NewDelete().
		Model(t).
		WherePK().
		Exec(ctx)
	if err != nil {
		return errors.New(errors.ErrCodeInternal, err.Error())
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		return errors.New(errors.ErrCodeInternal, err.Error())
	}

	if rows == 0 {
		return errors.New(errors.ErrCodeNotFound, "record not found")
	}

	return nil
}

func (r *userRepository) UpdateGeneric(ctx context.Context, logger *zap.Logger, t any) (int64, error) {
	if reflect.ValueOf(t).Kind().String() != "ptr" && reflect.ValueOf(t).Kind().String() != "struct" {
		return 0, errors.New(errors.ErrCodeInvalidArgument, "the argument must be a struct or pointer")
	}

	resp, err := r.db.NewUpdate().
		Model(t).
		OmitZero().
		WherePK().
		Exec(ctx)
	if err != nil {
		return 0, errors.New(errors.ErrCodeInternal, err.Error())
	}

	total, err := resp.RowsAffected()
	if err != nil {
		return 0, errors.New(errors.ErrCodeInternal, err.Error())
	}

	if total == 0 {
		return 0, errors.New(errors.ErrCodeNotFound, "record not found")
	}

	return total, nil
}
