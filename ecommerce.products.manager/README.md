# Ecommerce Products Manager

A REST API-based product management service for ecommerce applications, providing comprehensive product operations including CRUD operations, inventory management, and product search functionality.

## Features

- **Product Management**: Complete CRUD operations for product catalog
- **Inventory Management**: Stock tracking and update operations
- **Product Search**: Advanced search with filtering capabilities
- **REST API**: High-performance REST API with Gin framework
- **Database Integration**: PostgreSQL integration with Bun ORM
- **Structured Logging**: Comprehensive logging with Zap logger
- **Input Validation**: Robust request validation and error handling
- **API Documentation**: Swagger/OpenAPI documentation
- **Test Coverage**: Comprehensive test suite with high coverage
- **Metrics**: Built-in Prometheus metrics support
- **CORS Support**: Cross-Origin Resource Sharing enabled

## Architecture

The service follows a clean architecture pattern with the following layers:

```
├── cmd/                        # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/              # HTTP handlers (presentation layer)
│   │   ├── routes/                # Route definitions and middleware
│   │   │   └── middleware/            # Custom middleware
│   ├── config/                  # Configuration management
│   ├── errors/                  # Custom error handling
│   ├── models/                  # Data models and DTOs
│   ├── repository/              # Database layer
│   ├── service/                 # Business logic layer
│   ├── utils/                   # Utility functions and validation
│   └── mocks/                   # Test mocks
├── docs/                      # Swagger documentation
└── dockerfile                 # Docker configuration
```

## Prerequisites

- Go 1.23.2 or higher
- PostgreSQL database
- Swagger CLI (for API documentation)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Wembie/ecommerce-microservices.git
cd ecommerce-microservices
cd ecommerce.products.manager
```

2. Install dependencies:
```bash
go mod download
go mod tidy
```

3. Generate Swagger documentation (optional):
```bash
swag init -g cmd/main.go
```

## Configuration

Create a `.env` file with the following environment variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_INSECURE=false

# Connection Pool Settings
MAX_IDLE_CONN=10
MAX_OPEN_CONN=10

# Application Configuration
APP_NAME=products-manager
APP_PORT=8080

# Other Configuration
LOG_LEVEL=info
```

Initialize the environment variables:
```bash
set -a && source .env
```

## Usage

### Running the Service

Start the HTTP server:
```bash
go run cmd/main.go
```

The service will start listening on the configured port (default: 8080).

### API Endpoints

The service provides the following REST API endpoints:

#### Product Management
- `POST /api/products` - Create a new product
- `GET /api/products/{id}` - Get product by ID
- `PUT /api/products/{id}` - Update product information
- `DELETE /api/products/{id}` - Delete product
- `GET /api/products` - List products with pagination

#### Inventory Management
- `PATCH /api/products/{id}/stock` - Update product stock

#### Search & Filter
- `GET /api/products/search` - Search products with filters

#### API Documentation
- `GET /swagger/index.html` - Swagger UI documentation
- `GET /swagger/doc.json` - OpenAPI specification

### Example Usage

```bash
# Create a product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Wireless Headphones",
    "description": "High-quality wireless headphones with noise cancellation",
    "price": 299.99,
    "stock": 50
  }'

# Get a product
curl http://localhost:8080/api/products/{product-id}

# Search products
curl "http://localhost:8080/api/products/search?name=headphones&page=1&size=10"

# Update stock
curl -X PATCH http://localhost:8080/api/products/{product-id}/stock \
  -H "Content-Type: application/json" \
  -d '{"stock": 100}'
```

## Testing

### Running Tests

Execute all tests with coverage:
```bash
go test ./... -coverprofile=coverage.out
```

### Coverage Reports

Generate HTML coverage report:
```bash
go tool cover -html=coverage.out -o coverage.html
```

View function-level coverage:
```bash
go tool cover -func=coverage.out
```

### Test Structure

The project maintains **63.5%** overall test coverage across all packages:

**High Coverage Areas (80%+):**
- API handlers (97.6% - 100%)
- Middleware components (100%)
- Configuration management (100%)
- Error handling (100%)
- Service layer operations (100%)
- Utility functions (100%)

**Areas for Improvement:**
- Main application entry point (0%)
- Database repository layer (0%)
- Connection management (0%)
- Logging configuration (50% - 75%)
- Mock implementations (0% - testing utilities)

## Project Structure Details

### Core Components

**Configuration (`internal/config/`)**
- Environment-based configuration management
- Database connection setup with connection pooling
- Logging configuration

**Models (`internal/models/`)**
- Product entity definitions
- Request/Response DTOs
- Pagination models
- Database model mappings

**Repository (`internal/repository/`)**
- Generic database operations
- Product-specific database queries
- Database abstraction layer

**Service (`internal/service/`)**
- Business logic implementation
- Product validation and processing
- Inventory management logic

**API Layer (`internal/api/`)**
- HTTP handlers implementation
- Route definitions and middleware
- Request/response conversion
- Input validation

**Utils (`internal/utils/`)**
- Product validation utilities
- Helper functions
- Data transformation utilities

### Error Handling

The service implements a custom error system with:
- Structured error codes
- HTTP error mapping
- Comprehensive error logging
- Client-friendly error responses

### Security Features

- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries through Bun ORM
- **CORS Support**: Configurable cross-origin resource sharing
- **Request Timeout**: Built-in timeout middleware

### Performance Features

- **Connection Pooling**: Database connection pooling
- **Metrics**: Prometheus metrics integration
- **Structured Logging**: High-performance Zap logger
- **Pagination**: Efficient data pagination

## Dependencies

Key dependencies include:
- **Gin**: High-performance HTTP web framework
- **Bun**: PostgreSQL ORM
- **Zap**: Structured logging
- **UUID**: Unique identifier generation
- **Swagger**: API documentation generation
- **Testify**: Testing framework
- **CORS**: Cross-origin resource sharing middleware
- **Metrics**: Prometheus metrics support

## API Documentation

The service includes comprehensive Swagger/OpenAPI documentation:

1. **Interactive Documentation**: Available at `/swagger/index.html`
2. **API Specification**: JSON format at `/swagger/doc.json`
3. **Auto-generation**: Documentation generated from code annotations

To regenerate documentation:
```bash
swag init -g cmd/main.go
```

## Development

### Adding New Features

1. Define models in `internal/models/`
2. Implement repository methods in `internal/repository/`
3. Add business logic in `internal/service/`
4. Create API handlers in `internal/api/handlers/`
5. Add routes in `internal/api/routes/`
6. Add utility functions in `internal/utils/`
7. Write comprehensive tests
8. Update Swagger documentation

### Code Standards

- Follow Go best practices and conventions
- Maintain test coverage above 80%
- Use structured logging
- Implement proper error handling
- Add comprehensive documentation

## Docker

### Build and Run

1. Build the Docker image:
```bash
docker build -t ecommerce.products.manager:latest .
```

2. Run the container:
```bash
docker run -p 8080:8080 --name ecommerce-products-manager \
  -e DB_HOST=your_db_host \
  -e DB_USER=your_db_user \
  -e DB_PASSWORD=your_db_password \
  -e DB_NAME=your_db_name \
  -e APP_PORT=8080 \
  ecommerce.products.manager:latest
```

### Docker Compose

For development with PostgreSQL:
```yaml
version: '3.8'
services:
  products-manager:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USER=products_user
      - DB_PASSWORD=products_password
      - DB_NAME=products_db
      - APP_PORT=8080
    depends_on:
      - postgres
  
  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=products_user
      - POSTGRES_PASSWORD=products_password
      - POSTGRES_DB=products_db
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
```

## Monitoring

### Metrics

The service exposes Prometheus metrics at `/metrics` endpoint:
- HTTP request duration
- HTTP request count
- Database connection metrics
- Custom business metrics

### Logging

Structured logging with configurable levels:
- Request/response logging
- Error logging with stack traces
- Performance metrics
- Audit trails

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Ensure all tests pass
5. Maintain or improve test coverage
6. Follow the existing code structure
7. Add appropriate logging and error handling
8. Update documentation as needed
9. Submit a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

For issues and questions, please refer to the project's issue tracker or contact the development team.

## Changelog

### Version 1.0.0
- Initial release
- Product CRUD operations
- Inventory management
- Search functionality
- REST API with Swagger documentation
- PostgreSQL integration
- Comprehensive test suite