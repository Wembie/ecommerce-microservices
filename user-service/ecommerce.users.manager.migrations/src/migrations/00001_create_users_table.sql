-- ================================================================
-- Database
-- CREATE database user_service;
-- ================================================================

-- +goose Up
-- +goose StatementBegin

-- ================================================================
-- User Service - Users Table
-- ================================================================

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create dedicated schema
CREATE SCHEMA IF NOT EXISTS user_service;

-- ================================================================
-- 1. USERS
-- ================================================================
CREATE TABLE user_service.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT NULL,
    PRIMARY KEY (id)
);

-- ================================================================
-- INDEXES
-- ================================================================
CREATE INDEX idx_users_username ON user_service.users(username);
CREATE INDEX idx_users_email ON user_service.users(email);

-- ================================================================
-- COMMENTS AND DOCUMENTATION
-- ================================================================
COMMENT ON SCHEMA user_service IS 'Schema for user management';
COMMENT ON TABLE user_service.users IS 'Table storing user accounts with authentication info';
COMMENT ON COLUMN user_service.users.password_hash IS 'Hashed password for authentication';
COMMENT ON COLUMN user_service.users.created_at IS 'Timestamp when the user was created';
COMMENT ON COLUMN user_service.users.updated_at IS 'Timestamp when the user was last updated';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_service.users CASCADE;
DROP SCHEMA IF EXISTS user_service CASCADE;

-- +goose StatementEnd