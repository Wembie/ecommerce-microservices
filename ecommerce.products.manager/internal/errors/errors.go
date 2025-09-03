package errors

import (
	"errors"
	"net/http"
)

type ErrorAPI struct {
	StatusCode int    `json:"-"`
	Error      error  `json:"error"`
	Message    string `json:"message"`
}

func NewErrorAPI(statusCode int, err error, message string) *ErrorAPI {
	return &ErrorAPI{
		StatusCode: statusCode,
		Error:      err,
		Message:    message,
	}
}

var (
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDuplicateResource = errors.New("resource already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
)

var errorMappings = map[error]struct {
	Code    int
	Default string
}{
	ErrNotFound:          {http.StatusNotFound, "Resource not found"},
	ErrInvalidInput:      {http.StatusBadRequest, "Invalid input provided"},
	ErrDuplicateResource: {http.StatusConflict, "Resource already exists"},
	ErrUnauthorized:      {http.StatusUnauthorized, "Unauthorized access"},
	ErrForbidden:         {http.StatusForbidden, "Access forbidden"},
	ErrInternalServer:    {http.StatusInternalServerError, "Internal server error"},
}

func ValidateError(err error, customMessage ...string) *ErrorAPI {
	if err == nil {
		return nil
	}

	msg := ""
	if len(customMessage) > 0 && customMessage[0] != "" {
		msg = customMessage[0]
	}

	for key, meta := range errorMappings {
		if errors.Is(err, key) {
			if msg == "" {
				msg = meta.Default
			}
			return NewErrorAPI(meta.Code, err, msg)
		}
	}

	if msg == "" {
		msg = "Internal server error"
	}
	return NewErrorAPI(http.StatusInternalServerError, err, msg)
}
