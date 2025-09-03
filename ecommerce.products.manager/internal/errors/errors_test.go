package errors_test

import (
	"errors"
	"net/http"
	"testing"

	appErr "ecommerce.products.manager/internal/errors"
)

func TestValidateError(t *testing.T) {
	tests := []struct {
		name          string
		inputErr      error
		customMessage string
		wantCode      int
		wantMessage   string
	}{
		{
			name:        "ErrNotFound default",
			inputErr:    appErr.ErrNotFound,
			wantCode:    http.StatusNotFound,
			wantMessage: "Resource not found",
		},
		{
			name:          "ErrNotFound custom",
			inputErr:      appErr.ErrNotFound,
			customMessage: "Producto no encontrado",
			wantCode:      http.StatusNotFound,
			wantMessage:   "Producto no encontrado",
		},
		{
			name:        "ErrInvalidInput default",
			inputErr:    appErr.ErrInvalidInput,
			wantCode:    http.StatusBadRequest,
			wantMessage: "Invalid input provided",
		},
		{
			name:        "ErrDuplicateResource default",
			inputErr:    appErr.ErrDuplicateResource,
			wantCode:    http.StatusConflict,
			wantMessage: "Resource already exists",
		},
		{
			name:        "ErrUnauthorized default",
			inputErr:    appErr.ErrUnauthorized,
			wantCode:    http.StatusUnauthorized,
			wantMessage: "Unauthorized access",
		},
		{
			name:        "ErrForbidden default",
			inputErr:    appErr.ErrForbidden,
			wantCode:    http.StatusForbidden,
			wantMessage: "Access forbidden",
		},
		{
			name:        "Unknown error default",
			inputErr:    errors.New("something else"),
			wantCode:    http.StatusInternalServerError,
			wantMessage: "Internal server error",
		},
		{
			name:          "Unknown error custom",
			inputErr:      errors.New("random"),
			customMessage: "Algo raro pasó",
			wantCode:      http.StatusInternalServerError,
			wantMessage:   "Algo raro pasó",
		},
		{
			name:        "Nil error",
			inputErr:    nil,
			wantCode:    0,
			wantMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *appErr.ErrorAPI
			if tt.customMessage != "" {
				got = appErr.ValidateError(tt.inputErr, tt.customMessage)
			} else {
				got = appErr.ValidateError(tt.inputErr)
			}

			if tt.inputErr == nil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected non-nil, got nil")
			}

			if got.StatusCode != tt.wantCode {
				t.Errorf("expected code %d, got %d", tt.wantCode, got.StatusCode)
			}
			if got.Message != tt.wantMessage {
				t.Errorf("expected message %q, got %q", tt.wantMessage, got.Message)
			}
			if !errors.Is(got.Error, tt.inputErr) {
				t.Errorf("expected wrapped error %v, got %v", tt.inputErr, got.Error)
			}
		})
	}
}
