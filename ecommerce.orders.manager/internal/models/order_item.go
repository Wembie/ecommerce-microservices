package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OrderItem struct {
	bun.BaseModel `bun:"table:order_service.order_items"`

	ID        uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()" json:"id"`
	OrderID   uuid.UUID `bun:"order_id,notnull,type:uuid" json:"order_id"`
	ProductID uuid.UUID `bun:"product_id,notnull,type:uuid" json:"product_id"`
	Quantity  int       `bun:"quantity,notnull" json:"quantity"`
	Price     float64   `bun:"price,type:decimal(10,2),notnull" json:"price"`
}

type GetOrderItemsRequest struct {
	OrderID uuid.UUID `json:"order_id"`
}