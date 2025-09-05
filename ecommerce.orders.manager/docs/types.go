package docs

import (
	"ecommerce.orders.manager/internal/models"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message describing what went wrong"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Items        []models.OrderResponse `json:"items"`
	Page         int              `json:"page" example:"0"`
	Size         int              `json:"size" example:"50"`
	Total        int              `json:"total" example:"100"`
	Pages        int              `json:"pages" example:"2"`
	NextPage     *int             `json:"next_page,omitempty" example:"1"`
	PreviousPage *int             `json:"previous_page,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status" example:"healthy"`
}