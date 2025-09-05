# eCommerce Auth Manager

HTTP Authentication service for the eCommerce microservices architecture.

## Features

- HTTP REST API for authentication
- Proxies authentication to Users Manager via gRPC
- Token generation and validation endpoints

## API Endpoints

### POST /auth/token
OAuth2 password flow - exchange username/password for access token
- **Content-Type**: `application/x-www-form-urlencoded`
- **Body**: `username=user&password=pass`
- **Response**: `{"access_token": "jwt_token", "token_type": "bearer"}`

### POST /auth/validate  
Validates JWT token and returns user information
- **Content-Type**: `application/json`
- **Body**: `{"token": "jwt_token"}`
- **Response**: `{"valid": true, "user_id": "uuid", "username": "user", "email": "email"}`

### GET /health
Health check endpoint

## Configuration

Environment variables:

- `APP_PORT` - HTTP server port (default: 8080)
- `USER_MANAGER_HOST` - Users service gRPC address (default: dns:localhost:50051)

## Running

```bash
# Set environment variables
set -a && source .env

# Run the service
go run cmd/main.go
```

## Dependencies

- Users Manager service running on configured gRPC address
- Proto definitions from `github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo`