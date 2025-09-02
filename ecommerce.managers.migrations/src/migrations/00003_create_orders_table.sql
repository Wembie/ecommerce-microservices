-- ================================================================
-- Database
-- CREATE database ecommerce_db;
-- ================================================================

-- +goose Up
-- +goose StatementBegin

-- ================================================================
-- Order Service - Orders & Order Items Tables
-- ================================================================

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create dedicated schema
CREATE SCHEMA IF NOT EXISTS order_service;

-- ================================================================
-- 1. ORDERS
-- ================================================================
CREATE TABLE order_service.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total DECIMAL(10,2) NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT NULL,
    PRIMARY KEY (id)
);

-- ================================================================
-- 2. ORDER ITEMS
-- ================================================================
CREATE TABLE order_service.order_items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    order_id uuid NOT NULL,
    product_id uuid NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES order_service.orders(id) ON DELETE CASCADE
);

-- ================================================================
-- INDEXES
-- ================================================================
-- Orders
CREATE INDEX idx_orders_user_id ON order_service.orders(user_id);
CREATE INDEX idx_orders_status ON order_service.orders(status);
CREATE INDEX idx_orders_created_at ON order_service.orders(created_at);

-- Order Items
CREATE INDEX idx_order_items_order_id ON order_service.order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_service.order_items(product_id);

-- ================================================================
-- COMMENTS AND DOCUMENTATION
-- ================================================================
COMMENT ON SCHEMA order_service IS 'Schema for order management';
COMMENT ON TABLE order_service.orders IS 'Table storing customer orders';
COMMENT ON COLUMN order_service.orders.user_id IS 'Reference to the user who placed the order';
COMMENT ON COLUMN order_service.orders.status IS 'Current status of the order (pending, completed, etc.)';
COMMENT ON COLUMN order_service.orders.total IS 'Total cost of the order';
COMMENT ON COLUMN order_service.orders.created_at IS 'Timestamp when the order was created';
COMMENT ON COLUMN order_service.orders.updated_at IS 'Timestamp when the order was last updated';

COMMENT ON TABLE order_service.order_items IS 'Table storing individual items within an order';
COMMENT ON COLUMN order_service.order_items.order_id IS 'Reference to the related order';
COMMENT ON COLUMN order_service.order_items.product_id IS 'Reference to the purchased product';
COMMENT ON COLUMN order_service.order_items.quantity IS 'Quantity of the product in the order';
COMMENT ON COLUMN order_service.order_items.price IS 'Price of the product at the time of order';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS order_service.order_items CASCADE;
DROP TABLE IF EXISTS order_service.orders CASCADE;
DROP SCHEMA IF EXISTS order_service CASCADE;

-- +goose StatementEnd