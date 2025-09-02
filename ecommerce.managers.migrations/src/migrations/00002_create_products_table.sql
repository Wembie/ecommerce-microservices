-- ================================================================
-- Database
-- CREATE database product_service;
-- ================================================================

-- +goose Up
-- +goose StatementBegin

-- ================================================================
-- Product Service - Products Table
-- ================================================================

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create dedicated schema
CREATE SCHEMA IF NOT EXISTS product_service;

-- ================================================================
-- 1. PRODUCTS
-- ================================================================
CREATE TABLE product_service.products (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT NULL,
    PRIMARY KEY (id)
);

-- ================================================================
-- INDEXES
-- ================================================================
CREATE INDEX idx_products_name ON product_service.products(name);
CREATE INDEX idx_products_price ON product_service.products(price);

-- ================================================================
-- COMMENTS AND DOCUMENTATION
-- ================================================================
COMMENT ON SCHEMA product_service IS 'Schema for product management';
COMMENT ON TABLE product_service.products IS 'Table storing available products for sale';
COMMENT ON COLUMN product_service.products.name IS 'Name of the product';
COMMENT ON COLUMN product_service.products.description IS 'Detailed description of the product';
COMMENT ON COLUMN product_service.products.price IS 'Product price in decimal format';
COMMENT ON COLUMN product_service.products.stock IS 'Available stock units';
COMMENT ON COLUMN product_service.products.created_at IS 'Timestamp when the product was created';
COMMENT ON COLUMN product_service.products.updated_at IS 'Timestamp when the product was last updated';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS product_service.products CASCADE;
DROP SCHEMA IF EXISTS product_service CASCADE;

-- +goose StatementEnd
