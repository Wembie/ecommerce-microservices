# ecommerce.users.manager.lib.protos

User management service with gRPC and Protocol Buffers that provides authentication and user administration functionalities.

## Features

- User CRUD operations (create, read, update, delete)
- Username/password authentication
- JWT token validation
- gRPC API with Protocol Buffers

## Project Structure

```
v1/
│ ├── libgo/
│ │   ├── go.mod
│ │   ├── gosum
│ │   ├── users_grpc.pb.go
│ │   └── users.pb.go
│ ├── libpy/
│ │   ├── users_manager/
│ │   │   ├── __init__.py
│ │   │   ├── users_pb2_grpc.py
│ │   │   ├── users_pb2.py
│ │   │   └── users_pb2.pyi
│ │   ├── .gitignore
│ │   └── setup.py
│ ├── users_manager/
│ │   ├── buf.yaml
│ │   ├── users.proto
│ │   └── users.proto
│ ├── buf.gen.yaml
│ └── VERSION
├── .gitignore
└── README.md
```

## API

### UserService

| Method | Description |
|--------|-------------|
| `CreateUser` | Create new user |
| `GetUser` | Get user by ID |
| `UpdateUser` | Update user information |
| `DeleteUser` | Delete user |
| `AuthenticateUser` | Authenticate with credentials |
| `ValidateUser` | Validate JWT token |

### Main Messages

**CreateUserRequest**
- `username` (string)
- `email` (string) 
- `password` (string)

**AuthRequest**
- `username` (string)
- `password` (string)

**AuthResponse**
- `success` (bool)
- `token` (string, optional)
- `error_message` (string, optional)

## Usage

### Go Client
```go
conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := pb.NewUserServiceClient(conn)

user, _ := client.CreateUser(context.Background(), &pb.CreateUserRequest{
    Username: "Wembie",
    Email:    "wembie@example.com",
    Password: "secure_password",
})
```

### Python Client
```python
channel = grpc.insecure_channel('localhost:50051')
stub = users_pb2_grpc.UserServiceStub(channel)

response = stub.CreateUser(users_pb2.CreateUserRequest(
    username='Wembie',
    email='wembie@example.com',
    password='secure_password'
))
```

## Generate Protocol Buffers

```bash
# Navigate to v1 directory
cd v1

# Generate code
buf generate
```