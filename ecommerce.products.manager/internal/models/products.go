package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product_service.products"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	Name        string     `bun:"name,notnull" json:"name"`
	Description *string    `bun:"description,nullzero" json:"description,omitempty"`
	Price       float64    `bun:"price,type:decimal(10,2),notnull" json:"price"`
	Stock       int        `bun:"stock,notnull" json:"stock"`
	CreatedAt   time.Time  `bun:"created_at,default:now()" json:"created_at"`
	UpdatedAt   *time.Time `bun:"updated_at,nullzero" json:"updated_at,omitempty"`
}

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
}

type GetProductRequest struct {
	ID uuid.UUID `json:"id"`
}

type UpdateProductRequest struct {
	bun.BaseModel `bun:"table:product_service.products"`

	ID          uuid.UUID  `bun:"id,pk,type:uuid" json:"-"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	Stock       *int       `json:"stock,omitempty"`
	UpdatedAt   time.Time  `bun:"updated_at,type:timestamptz" json:"-"`
}

type DeleteProductRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteProductResponse struct {
	Success bool `json:"success"`
}

type UpdateStockRequest struct {
	bun.BaseModel `bun:"table:product_service.products"`

	ID    uuid.UUID	`bun:"id,pk,type:uuid" json:"-"`
	Stock int       `json:"stock"`
}

type SearchProductsRequest struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	Stock       *int      `json:"stock,omitempty"`
	Page        int       `json:"page"`
	Size        int       `json:"size"`
}