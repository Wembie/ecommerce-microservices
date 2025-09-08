# 🛒 E-commerce Microservices

A comprehensive e-commerce platform built with Go microservices architecture, featuring user management, product catalog, order processing, and authentication services.

## 🏗️ Architecture Overview

This project consists of 5 main microservices:

- **🔐 Authentication Service** (`ecommerce.auth`) - JWT-based authentication and authorization
- **👥 User Manager** (`ecommerce.users.manager`) - User management via gRPC
- **📦 Product Manager** (`ecommerce.products.manager`) - Product catalog management with REST API
- **🛍️ Order Manager** (`ecommerce.orders.manager`) - Order processing and management
- **🗄️ Database Migrations** (`ecommerce.managers.migrations`) - Database schema management

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.23.2+ (for development)
- PostgreSQL 15 (handled by Docker)

### 🐳 Running with Docker Compose

1. **Clone the repository**
   ```bash
   git clone https://github.com/Wembie/ecommerce-microservices.git
   cd ecommerce-microservices
   ```

2. **Start all services**
   ```bash
   docker-compose up --build
   ```

3. **Check service health**
   ```bash
   docker-compose ps
   ```

4. **Stop all services**
   ```bash
   docker-compose down
   ```

5. **Stop and remove volumes** (⚠️ This will delete database data)
   ```bash
   docker-compose down -v
   ```

### 📋 Service Endpoints

| Service | Port | Type | Documentation |
|---------|------|------|---------------|
| Authentication | 8080 | REST | http://localhost:8080/swagger/index.html |
| Products Manager | 8081 | REST | http://localhost:8081/swagger/index.html |
| Orders Manager | 8082 | REST | http://localhost:8082/swagger/index.html |
| User Manager | 50051 | gRPC | - |
| PostgreSQL | 5434 | Database | - |

## 🔧 Development Setup

### Building Individual Services

Each service can be built and run independently:

```bash
# User Manager (gRPC Service)
cd ecommerce.users.manager
go mod tidy
go run cmd/main.go

# Product Manager (REST API)
cd ecommerce.products.manager
go mod tidy
go run cmd/main.go

# Order Manager (REST API)
cd ecommerce.orders.manager
go mod tidy
go run cmd/main.go

# Authentication Service (REST API)
cd ecommerce.auth
go mod tidy
go run cmd/main.go
```

### Running Tests

```bash
# Run tests for all services
cd ecommerce.users.manager && go test ./...
cd ecommerce.products.manager && go test ./...
cd ecommerce.orders.manager && go test ./...
cd ecommerce.auth && go test ./...
```

### Generating API Documentation

Each REST service has Swagger documentation:

```bash
# Generate docs for Products Manager
cd ecommerce.products.manager
./scripts/generate-docs.sh  # or generate-docs.bat on Windows

# Generate docs for Orders Manager
cd ecommerce.orders.manager
./scripts/generate-docs.sh

# Generate docs for Auth Service
cd ecommerce.auth
./scripts/generate-docs.sh
```

## 🗄️ Database

### Schema Management

Database migrations are handled automatically by the `migrations` service during startup.

### Connection Details

- **Host**: localhost
- **Port**: 5434
- **Database**: ecommerce_db
- **User**: app_user
- **Password**: app_password

### Manual Database Access

```bash
docker exec -it ecommerce-db psql -U app_user -d ecommerce_db
```

## 🔐 Environment Variables

### Common Database Variables
```env
DB_HOST=db
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=app_password
DB_NAME=ecommerce_db
DB_INSECURE=true
MAX_IDLE_CONN=10
MAX_OPEN_CONN=10
```

### Service-Specific Variables

#### Authentication Service
```env
APP_PORT=8080
USER_MANAGER_HOST=user-manager:50051
```

#### User Manager
```env
APP_NAME=users-manager
APP_PORT=50051
JWT_SECRET=your_jwt_secret
```

#### Product Manager
```env
APP_NAME=products-manager
APP_PORT=8081
```

#### Order Manager
```env
APP_NAME=orders-manager
APP_PORT=8082
USER_MANAGER_HOST=user-manager:50051
PRODUCT_MANAGER_HOST=http://product-manager:8081
```

## 🏢 Service Details

### 🔐 Authentication Service
- **Technology**: Gin (REST API)
- **Port**: 8080
- **Features**: JWT token generation, user authentication
- **Dependencies**: User Manager (gRPC)

### 👥 User Manager
- **Technology**: gRPC
- **Port**: 50051
- **Features**: User CRUD operations, password management
- **Database**: Direct PostgreSQL connection

### 📦 Product Manager
- **Technology**: Gin (REST API)
- **Port**: 8081
- **Features**: Product catalog, inventory management, search, pagination
- **Database**: Direct PostgreSQL connection
- **API Documentation**: Swagger/OpenAPI

### 🛍️ Order Manager
- **Technology**: Gin (REST API)
- **Port**: 8082
- **Features**: Order processing, order history, inventory validation
- **Dependencies**: User Manager (gRPC), Product Manager (HTTP)
- **API Documentation**: Swagger/OpenAPI

### 🗄️ Database Migrations
- **Technology**: Go + PostgreSQL
- **Features**: Automatic schema setup and migrations
- **Runs**: On container startup before other services

## 📊 Monitoring & Health Checks

All services include:
- Health check endpoints
- Structured logging with Zap
- Metrics collection with Prometheus (where applicable)
- CORS support for web applications

## 🔗 Service Communication

```
┌─────────────────┐    HTTP     ┌─────────────────┐
│  Authentication │────────────▶│   User Manager  │
│    Service      │    gRPC     │     (gRPC)      │
│    (REST)       │             │                 │
└─────────────────┘             └─────────────────┘
                                          │
                                       gRPC
                                          ▼
┌─────────────────┐    HTTP     ┌─────────────────┐
│  Order Manager  │────────────▶│ Product Manager │
│    (REST)       │             │     (REST)      │
└─────────────────┘             └─────────────────┘
         │                               │
         └──────────── PostgreSQL ──────┘
```

## 🚫 Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 5434, 8080, 8081, 8082, 50051 are available
2. **Database connection errors**: Wait for the database health check to pass
3. **Service startup order**: Docker Compose handles dependencies automatically

### Logs

```bash
# View all service logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f user-manager
docker-compose logs -f product-manager
docker-compose logs -f auth
```

## 🧪 Testing the API

### Authentication Flow
1. Use User Manager to create a user (gRPC)
2. Use Authentication Service to login and get JWT token
3. Use the JWT token to access protected endpoints in other services

### Example API Calls

```bash
# Health check
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health

# Access Swagger UI
open http://localhost:8080/swagger/index.html
open http://localhost:8081/swagger/index.html
open http://localhost:8082/swagger/index.html
```

## 📮 Postman Collections

### 🔗 Shared Collections

For easy API testing, Postman collections are available:

#### HTTP REST APIs
- **Shared Link**: [Open in Postman](https://postman.co/workspace/My-Workspace~5d1ca1e6-a626-44eb-8391-51a7b92f8149/collection/68b4764565db39b34814f374?action=share&creator=37797695)
- **Local JSON File**: `Ecommerce HTTP.postman_collection.json` (Import this file into Postman)

#### gRPC APIs  
- **Shared Link**: [Open in Postman](https://postman.co/workspace/My-Workspace~5d1ca1e6-a626-44eb-8391-51a7b92f8149/collection/68b4764565db39b34814f374?action=share&creator=37797695)  
- **Note**: gRPC collection couldn’t be exported as JSON, please use the shared link above

### 📥 How to Use Postman Collections

1. **Using Shared Links**: Click the links above to access the collections directly in Postman web
2. **Using JSON File**: 
   ```bash
   # Import the HTTP collection
   # 1. Open Postman
   # 2. Click "Import" 
   # 3. Select "Ecommerce HTTP.postman_collection.json"
   ```
3. **Environment Setup**: Make sure to set the base URLs:
   - `auth_url`: http://localhost:8080
   - `products_url`: http://localhost:8081  
   - `orders_url`: http://localhost:8082
   - `users_grpc_url`: localhost:50051

## 📈 Current Status

- ✅ User Management (gRPC)
- ✅ Authentication Service (REST)
- ✅ Product Management (REST)
- ✅ Order Management (REST)
- ✅ Database Migrations
- ✅ Docker Compose Setup
- ✅ API Documentation (Swagger)
- ✅ Comprehensive Testing Suite

## 🔄 Development Workflow

1. Make changes to service code
2. Run tests: `go test ./...`
3. Update Swagger docs if API changes: `./scripts/generate-docs.sh`
4. Rebuild and test with Docker Compose: `docker-compose up --build`
5. Commit changes to feature branch
6. Create pull request to `main` branch

## 📝 License

This project is licensed under the terms specified in the LICENSE file.
