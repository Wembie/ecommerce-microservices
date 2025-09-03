package utils

import (
	"net/http"

	"ecommerce.products.manager/internal/errors"

	"go.uber.org/zap"
)

func ValidateProductFields(logger *zap.Logger, name *string, price *float64, stock *int, description *string, requireName bool) *errors.ErrorAPI {
	if name != nil {
		if *name == "" {
			logger.Warn("Product name cannot be empty")
			return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Product name cannot be empty")
		}
		if len(*name) > 255 {
			logger.Warn("Product name too long")
			return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Product name cannot exceed 255 characters")
		}
	} else if requireName {
		logger.Warn("Product name is required")
		return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Product name is required and cannot be empty")
	}

	if price != nil && *price <= 0 {
		logger.Warn("Invalid product price")
		return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Product price must be greater than 0")
	}

	if stock != nil && *stock < 0 {
		logger.Warn("Invalid product stock")
		return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Product stock cannot be negative")
	}

	if description != nil && len(*description) > 1000 {
		logger.Warn("Product description too long")
		return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Product description cannot exceed 1000 characters")
	}

	return nil
}