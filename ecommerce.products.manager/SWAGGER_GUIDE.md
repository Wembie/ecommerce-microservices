# 📚 Swagger Documentation Guide - Products Manager API

> **Created by:** Wembie  
> **Date:** September 2025  
> **API Version:** 1.0.0

## 🎯 Summary of Changes Made

### ✅ **Implemented Improvements:**

1. **🔄 ListProducts Elimination**: Removed the redundant `GET /products` endpoint
2. **🎯 REST-compliant SearchProducts**: **CHANGED TO GET** - Now uses `GET /products/search` with query parameters (following REST standards)
3. **📝 Enhanced Query Parameters**: Filters using query params for better cacheability and usability
4. **🔧 Updated Pagination**: Changed to 0-based pagination with defaults `page=0, size=50`
5. **✅ Comprehensive Validations**: Added robust validations across all endpoints
6. **📖 Complete Swagger Documentation**: Professional configuration with author information

### 🔄 **Why GET instead of POST for searches?**

**REST Principles:**
- ✅ **GET** is idempotent (multiple calls = same result)
- ✅ **GET** is cacheable by browsers/proxies
- ✅ **GET** allows bookmarking of searches
- ✅ **GET** is semantically correct for read operations

---

## 🚀 How to Generate Swagger Documentation

### **Option 1: Automated Scripts (Recommended)**

#### Windows:
```batch
# Run from project root
scripts\generate-docs.bat
```

#### Linux/macOS:
```bash
# Run from project root
./scripts/generate-docs.sh
```

### **Option 2: Manual Command**

```bash
# 1. Install swag CLI (first time only)
go install github.com/swaggo/swag/cmd/swag@latest

# 2. Clean previous docs (optional)
rm -rf docs/

# 3. Generate documentation
swag init -g cmd/main.go -o docs --parseDependency --parseInternal
```

### **Parameters Explained:**

- `-g cmd/main.go`: Main application file
- `-o docs`: Output directory for generated files
- `--parseDependency`: Parses external dependencies
- `--parseInternal`: Parses internal project packages

---

## 📋 API Structure

### **🔍 Main Endpoint: SearchProducts**

**URL:** `GET /products/search`

**Query Parameters:**
```
?name=iPhone              // Optional: filter by name (partial search)
&description=smartphone   // Optional: filter by description (partial search)  
&price=999.99            // Optional: filter by exact price
&stock=10                // Optional: filter by minimum stock
&page=0                  // Page (0-based, default: 0)
&size=50                 // Items per page (default: 50)
```

**Complete Example:**
```
GET /products/search?name=iPhone&price=999.99&page=0&size=20
```

**Response:**
```json
{
  "items": [
    {
      "id": "uuid",
      "name": "iPhone 15",
      "description": "Latest smartphone",
      "price": 999.99,
      "stock": 25,
      "created_at": "2025-09-02T10:00:00Z",
      "updated_at": "2025-09-02T12:00:00Z"
    }
  ],
  "page": 0,
  "size": 50,
  "total": 100,
  "pages": 2,
  "next_page": 1,
  "previous_page": null
}
```

### **✨ Other Endpoints:**

- `POST /products` - Create product
- `GET /products/{id}` - Get product by ID
- `PUT /products/{id}` - Update product
- `DELETE /products/{id}` - Delete product
- `PUT /products/{id}/stock` - Update stock

---

## 🔧 Configuration and Usage

### **1. Start the Application**

```bash
# Development
go run cmd/main.go

# Production
go build -o app cmd/main.go
./app
```

### **2. Access Documentation**

Once the application is running:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON Spec**: http://localhost:8080/swagger/doc.json
- **Health Check**: http://localhost:8080/health

### **3. Environment Variables**

```bash
# Optional: Configure host and version
export SWAGGER_HOST="localhost:8080"
export APP_VERSION="1.0.0"
export SWAGGER_BASE_PATH=""  # Empty by default
```

---

## 📖 Generated Files

After running `swag init`, these files are generated:

```
docs/
├── docs.go        # Go code with embedded specification
├── swagger.json   # OpenAPI 3.0 specification in JSON
├── swagger.yaml   # OpenAPI 3.0 specification in YAML
└── types.go       # Data types for Swagger
```

---

## 🎨 Customization

### **Modify API Information**

Edit the comments in `internal/docs/swagger.go`:

```go
// @title Products Manager API - Wembie
// @version 1.0.0
// @description Your custom description here
// @contact.name Your Name
// @contact.email your.email@example.com
```

### **Add New Endpoints**

1. Add Swagger comments in your handler:
```go
// NewEndpoint godoc
// @Summary Short description
// @Description Long description of the endpoint
// @Tags Products
// @Accept json
// @Produce json
// @Param request body ModelRequest true "Parameter description"
// @Success 200 {object} ModelResponse
// @Failure 400 {object} docs.ErrorResponse
// @Router /new-endpoint [post]
func (h *Handler) NewEndpoint(c *gin.Context) {
    // Implementation
}
```

2. Regenerate documentation:
```bash
swag init -g cmd/main.go -o docs --parseDependency --parseInternal
```

---

## ❓ Troubleshooting

### **Error: "cannot find type definition"**
- **Solution**: Verify that all types are correctly imported in `docs/types.go`

### **Error: "failed to get package name"**
- **Solution**: Make sure you're running the command from the project root

### **Swagger UI doesn't load**
- **Solution**: Verify that the application is running and the `/swagger/*any` route is configured

### **Missing fields in documentation**
- **Solution**: Add JSON tags and Swagger comments to your structs:
```go
type Product struct {
    ID   string `json:"id" example:"uuid"`
    Name string `json:"name" example:"iPhone 15"`
}
```

---

## 🏗 Senior Best Practices

### **1. Robust Validations**
```go
// Validate inputs
if request.Name == "" {
    return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Name is required")
}

if len(request.Name) > 255 {
    return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Name too long")
}
```

### **2. Structured Logging**
```go
logger.Info("Creating Product", zap.Any("request", request))
logger.Warn("Invalid input", zap.String("field", "name"), zap.String("value", request.Name))
```

### **3. Consistent Responses**
```go
// Use PaginatedResponse model for consistency
response := models.NewPaginatedResponse(items, page, size, total)
c.JSON(http.StatusOK, response)
```

### **4. Error Handling**
```go
// Structured errors with appropriate HTTP codes
if err != nil {
    logger.Error("Database error", zap.Error(err))
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
    return
}
```

---

## 🎯 Usage Example with cURL

```bash
# Search products (new GET endpoint - REST correct)
curl -X GET "http://localhost:8080/products/search?name=iPhone&page=0&size=10"

# Search with multiple filters
curl -X GET "http://localhost:8080/products/search?name=iPhone&price=999.99&stock=5&page=0&size=20"

# Search with pagination only (list all)
curl -X GET "http://localhost:8080/products/search?page=0&size=50"

# Create product
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest smartphone",
    "price": 999.99,
    "stock": 50
  }'
```

---

**Complete and professional documentation ready to use! 🚀**

> For any questions or improvements, contact **Wembie** - The architect of this API 😎