# ecommerce-microservices
Backend technical test project implementing a microservices-based e-commerce system in Go.   Includes three services (User, Product, Order) using gRPC and REST, PostgreSQL with migrations,   Dockerized setup with docker-compose, and documented APIs (Swagger &amp; Protobuf).  

docker-compose up --build
docker compose up --build
docker-compose down -v

### Current Test Coverage

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
