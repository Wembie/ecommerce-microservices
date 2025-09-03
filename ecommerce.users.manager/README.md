# Ecommerce Users Manager

A gRPC-based user management service for ecommerce applications, providing comprehensive user operations including authentication, validation, and CRUD operations.

## Features

- **User Management**: Complete CRUD operations for user accounts
- **Authentication**: JWT-based user authentication system
- **Token Validation**: Secure token validation for authorized access
- **gRPC API**: High-performance gRPC service interface
- **Database Integration**: PostgreSQL integration with Bun ORM
- **Structured Logging**: Comprehensive logging with Zap logger and trace ID support
- **Input Validation**: Robust request validation and error handling
- **Test Coverage**: Comprehensive test suite with 65.5% overall coverage

## Architecture

The service follows a clean architecture pattern with the following layers:

```
├── cmd/                    # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── errors/            # Custom error handling
│   ├── models/            # Data models and DTOs
│   ├── repository/        # Database layer
│   ├── service/           # Business logic layer
│   ├── transport/         # gRPC handlers (presentation layer)
│   ├── utils/             # Utility functions and converters
│   └── mocks/             # Test mocks
```

## Prerequisites

- Go 1.23.2 or higher
- PostgreSQL database
- Protocol Buffers compiler (protoc)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Wembie/ecommerce-microservices.git
cd user-service
cd ecommerce.users.manager
```

2. Install dependencies:
```bash
go mod download
go mod tidy
```

3. Install the custom protobuf library:
```bash
go get -v github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo
```
o
```bash
go get github.com/Wembie/ecommerce.users.manager.lib.protos/v1/libgo@latest
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
APP_NAME=users-manager
APP_PORT=50051

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
```

Initialize the environment variables:
```bash
set -a && source .env
```

## Usage

### Running the Service

Start the gRPC server:
```bash
go run cmd/main.go
```

The service will start listening on the configured port (default: 50051).

### API Operations

The service provides the following gRPC operations:

#### User Management
- `CreateUser`: Create a new user account
- `GetUser`: Retrieve user information by ID
- `UpdateUser`: Update user profile information
- `DeleteUser`: Remove user account

#### Authentication
- `AuthenticateUser`: Authenticate user credentials and generate JWT token
- `ValidateUser`: Validate JWT token and retrieve user information

### Example Usage

```go
// Example gRPC client usage
conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb.NewUserServiceClient(conn)

// Create a user
createReq := &pb.CreateUserRequest{
    Username: "Wembie",
    Email:    "wembie@example.com",
    Password: "secure_password",
}

user, err := client.CreateUser(context.Background(), createReq)
if err != nil {
    log.Fatal(err)
}

// Authenticate user
authReq := &pb.AuthRequest{
    Username: "Wembie",
    Password: "secure_password",
}

authResp, err := client.AuthenticateUser(context.Background(), authReq)
if err != nil {
    log.Fatal(err)
}
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

### Current Test Coverage

The project maintains **65.5%** overall test coverage across all packages:

**High Coverage Areas (80%+):**
- Configuration management (100%)
- Error handling (100%) 
- User service operations (81.8% - 100%)
- Transport handlers (87.5% - 100%)
- Utility functions (95.2% - 100%)

**Areas for Improvement:**
- Main application entry point (0%)
- Database repository layer (0%)
- Mock implementations (0% - testing utilities)

## Project Structure Details

### Core Components

**Configuration (`internal/config/`)**
- Environment-based configuration management
- Database connection setup
- Logging configuration with trace ID support

**Models (`internal/models/`)**
- User entity definitions
- Request/Response DTOs
- Database model mappings

**Repository (`internal/repository/`)**
- Generic database operations
- User-specific database queries
- Database abstraction layer

**Service (`internal/service/`)**
- Business logic implementation
- Password hashing and JWT token management
- User validation and authentication

**Transport (`internal/transport/`)**
- gRPC handler implementations
- Request/response conversion
- Protocol buffer integration

**Utils (`internal/utils/`)**
- Protocol buffer converters
- Helper functions
- Data transformation utilities

### Error Handling

The service implements a custom error system with:
- Structured error codes
- gRPC error mapping
- Comprehensive error logging

### Security Features

- **Password Hashing**: Secure password storage
- **JWT Authentication**: Token-based authentication
- **Input Validation**: Email format validation and required field checks
- **SQL Injection Protection**: Parameterized queries through Bun ORM

## Dependencies

Key dependencies include:
- **gRPC**: High-performance RPC framework
- **Bun**: PostgreSQL ORM
- **Zap**: Structured logging
- **UUID**: Unique identifier generation
- **GoValidator**: Input validation
- **Testify**: Testing framework

## Development

### Adding New Features

1. Define models in `internal/models/`
2. Implement repository methods in `internal/repository/`
3. Add business logic in `internal/service/`
4. Create transport handlers in `internal/transport/`
5. Add utility converters in `internal/utils/`
6. Write comprehensive tests

## Docker

1. Build
    ```bash
    docker build -t ecommerce.users.manager:latest .
    ```
2. Run
    ```bash
    docker run -p 50051:50051 --name ecommerce-users-manager ecommerce.users.manager:latest
    ```

### Contributing

1. Ensure all tests pass
2. Maintain or improve test coverage
3. Follow the existing code structure
4. Add appropriate logging and error handling
5. Update documentation as needed

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please refer to the project's issue tracker or contact the development team.