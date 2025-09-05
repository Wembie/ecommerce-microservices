package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:order_service.orders"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `bun:"user_id,notnull,type:uuid" json:"user_id"`
	Status    string     `bun:"status,notnull" json:"status"`
	Total     float64    `bun:"total,type:decimal(10,2),notnull" json:"total"`
	CreatedAt time.Time  `bun:"created_at,default:now()" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero" json:"updated_at,omitempty"`
}

type CreateOrderRequest struct {
	UserID uuid.UUID         `json:"-"`
	Items  []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type GetOrderRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetOrdersByUserRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Page   int       `json:"page"`
	Size   int       `json:"size"`
}

type UpdateOrderStatusRequest struct {
	bun.BaseModel `bun:"table:order_service.orders"`

	ID        uuid.UUID `bun:"id,pk,type:uuid" json:"-"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamptz" json:"-"`
}

type OrderResponse struct {
	Order *Order       `json:"order"`
	Items []*OrderItem `json:"items"`
}