package docs

import (
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message describing what went wrong"`
}