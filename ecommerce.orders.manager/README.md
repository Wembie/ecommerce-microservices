# Order Service - eCommerce Microservices

## Overview

The Order Service is a microservice built in Go that handles order management operations including creating orders, updating order status, and retrieving order information. It communicates with the User Service via gRPC for user validation and with the Product Service via HTTP for product information and stock management.

## Features

- **Order Management**: Create, retrieve, update, and manage orders
- **Authentication**: JWT-based authentication via User Service
- **Product Integration**: Real-time product validation and stock management
- **Database**: PostgreSQL with Bun ORM
- **Logging**: Structured logging with Zap
- **Metrics**: Prometheus metrics integration
- **Documentation**: Swagger API documentation

## Architecture

- **HTTP REST API** for client communication
- **gRPC Client** for User Service integration
- **HTTP Client** for Product Service integration
- **PostgreSQL** for data persistence
- **JWT Middleware** for authentication
- **Clean Architecture** with separated layers

## API Endpoints

### Orders
- `POST /orders` - Create a new order
- `GET /orders/{id}` - Get order by ID
- `GET /orders/user/{user_id}` - Get orders by user ID (paginated)
- `PUT /orders/{id}/status` - Update order status
- `GET /orders/{id}/items` - Get order items

### Health & Monitoring
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics
- `GET /swagger/*` - Swagger documentation

## Configuration

Environment variables:
- `DB_HOST` - Database host
- `DB_PORT` - Database port  
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `APP_PORT` - Application port
- `APP_NAME` - Application name
- `USER_MANAGER_HOST` - User service gRPC address (default: dns:localhost:50051)
- `PRODUCT_MANAGER_HOST` - Product service HTTP URL (default: http://localhost:8080)

## Dependencies

- User Service (gRPC on port 50051)
- Product Service (HTTP on port 8080)
- PostgreSQL database

## Running the Service

```bash
# Set environment variables
set -a && source .env

# Run the service
go run cmd/main.go
```

## Database Schema

The service uses the following tables in the `order_service` schema:

- `orders` - Main order information
- `order_items` - Order line items

## Authentication

All endpoints require JWT authentication via the `Authorization: Bearer <token>` header. The token is validated through the User Service.

## Error Handling

The service returns structured error responses with appropriate HTTP status codes:
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (invalid/missing token)
- `404` - Not Found (resource not found)
- `500` - Internal Server Error