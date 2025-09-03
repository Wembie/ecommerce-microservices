# 📚 Guía de Documentación Swagger - Products Manager API

> **Creado por:** Wembie  
> **Fecha:** Septiembre 2025  
> **API Version:** 1.0.0

## 🎯 Resumen de los Cambios Realizados

### ✅ **Mejoras Implementadas:**

1. **🔄 Eliminación de ListProducts**: Se removió el endpoint redundante `GET /products`
2. **🎯 SearchProducts REST-compliant**: **CAMBIADO A GET** - Ahora usa `GET /products/search` con query parameters (siguiendo estándares REST)
3. **📝 Query Parameters mejorados**: Filtros usando query params para mejor cacheabilidad y usabilidad
4. **🔧 Paginación actualizada**: Cambio a paginación base-0 con defaults `page=0, size=50`
5. **✅ Validaciones comprehensivas**: Agregadas validaciones robustas en todos los endpoints
6. **📖 Documentación Swagger completa**: Configuración profesional con información de autor

### 🔄 **¿Por qué GET en lugar de POST para búsquedas?**

**Principios REST:**
- ✅ **GET** es idempotente (múltiples llamadas = mismo resultado)
- ✅ **GET** es cacheable por navegadores/proxies
- ✅ **GET** permite bookmarking de búsquedas
- ✅ **GET** es el método semánticamente correcto para operaciones de lectura

---

## 🚀 Cómo Generar la Documentación Swagger

### **Opción 1: Scripts Automatizados (Recomendado)**

#### Windows:
```batch
# Ejecutar desde la raíz del proyecto
scripts\generate-docs.bat
```

#### Linux/macOS:
```bash
# Ejecutar desde la raíz del proyecto
./scripts/generate-docs.sh
```

### **Opción 2: Comando Manual**

```bash
# 1. Instalar swag CLI (solo primera vez)
go install github.com/swaggo/swag/cmd/swag@latest

# 2. Limpiar docs anteriores (opcional)
rm -rf docs/

# 3. Generar documentación
swag init -g cmd/main.go -o docs --parseDependency --parseInternal
```

### **Parámetros Explicados:**

- `-g cmd/main.go`: Archivo principal de la aplicación
- `-o docs`: Directorio de salida para los archivos generados
- `--parseDependency`: Analiza dependencias externas
- `--parseInternal`: Analiza paquetes internos del proyecto

---

## 📋 Estructura de la API

### **🔍 Endpoint Principal: SearchProducts**

**URL:** `GET /products/search`

**Query Parameters:**
```
?name=iPhone              // Opcional: filtrar por nombre (búsqueda parcial)
&description=smartphone   // Opcional: filtrar por descripción (búsqueda parcial)  
&price=999.99            // Opcional: filtrar por precio exacta
&stock=10                // Opcional: filtrar por stock mínimo
&page=0                  // Página (base-0, default: 0)
&size=50                 // Elementos por página (default: 50)
```

**Ejemplo completo:**
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

### **✨ Otros Endpoints:**

- `POST /products` - Crear producto
- `GET /products/{id}` - Obtener producto por ID
- `PUT /products/{id}` - Actualizar producto
- `DELETE /products/{id}` - Eliminar producto
- `PUT /products/{id}/stock` - Actualizar stock

---

## 🔧 Configuración y Uso

### **1. Iniciar la Aplicación**

```bash
# Desarrollo
go run cmd/main.go

# Producción
go build -o app cmd/main.go
./app
```

### **2. Acceder a la Documentación**

Una vez que la aplicación esté ejecutándose:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON Spec**: http://localhost:8080/swagger/doc.json
- **Health Check**: http://localhost:8080/health

### **3. Variables de Entorno**

```bash
# Opcional: Configurar host y versión
export SWAGGER_HOST="localhost:8080"
export APP_VERSION="1.0.0"
export SWAGGER_BASE_PATH=""  # Vacío por defecto
```

---

## 📖 Archivos Generados

Después de ejecutar `swag init`, se generan estos archivos:

```
docs/
├── docs.go        # Código Go con la especificación embebida
├── swagger.json   # Especificación OpenAPI 3.0 en JSON
├── swagger.yaml   # Especificación OpenAPI 3.0 en YAML
└── types.go       # Tipos de datos para Swagger
```

---

## 🎨 Personalización

### **Modificar Información de la API**

Edita los comentarios en `internal/docs/swagger.go`:

```go
// @title Products Manager API - Wembie
// @version 1.0.0
// @description Tu descripción personalizada aquí
// @contact.name Tu Nombre
// @contact.email tu.email@example.com
```

### **Agregar Nuevos Endpoints**

1. Agrega comentarios Swagger en tu handler:
```go
// NuevoEndpoint godoc
// @Summary Descripción corta
// @Description Descripción larga del endpoint
// @Tags Products
// @Accept json
// @Produce json
// @Param request body ModelRequest true "Descripción del parámetro"
// @Success 200 {object} ModelResponse
// @Failure 400 {object} docs.ErrorResponse
// @Router /nuevo-endpoint [post]
func (h *Handler) NuevoEndpoint(c *gin.Context) {
    // Implementación
}
```

2. Regenera la documentación:
```bash
swag init -g cmd/main.go -o docs --parseDependency --parseInternal
```

---

## ❓ Troubleshooting

### **Error: "cannot find type definition"**
- **Solución**: Verifica que todos los tipos estén importados correctamente en `docs/types.go`

### **Error: "failed to get package name"**
- **Solución**: Asegúrate de estar ejecutando el comando desde la raíz del proyecto

### **Swagger UI no carga**
- **Solución**: Verifica que la aplicación esté corriendo y que la ruta `/swagger/*any` esté configurada

### **Campos faltantes en la documentación**
- **Solución**: Agrega tags JSON y comentarios Swagger a tus structs:
```go
type Product struct {
    ID   string `json:"id" example:"uuid"`
    Name string `json:"name" example:"iPhone 15"`
}
```

---

## 🏗 Best Practices Senior

### **1. Validaciones Robustas**
```go
// Validar inputs
if request.Name == "" {
    return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Name is required")
}

if len(request.Name) > 255 {
    return errors.NewErrorAPI(http.StatusBadRequest, errors.ErrInvalidInput, "Name too long")
}
```

### **2. Logging Estructurado**
```go
logger.Info("Creating Product", zap.Any("request", request))
logger.Warn("Invalid input", zap.String("field", "name"), zap.String("value", request.Name))
```

### **3. Responses Consistentes**
```go
// Usar el modelo PaginatedResponse para consistencia
response := models.NewPaginatedResponse(items, page, size, total)
c.JSON(http.StatusOK, response)
```

### **4. Error Handling**
```go
// Errors estructurados con códigos HTTP apropiados
if err != nil {
    logger.Error("Database error", zap.Error(err))
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
    return
}
```

---

## 🎯 Ejemplo de Uso con cURL

```bash
# Buscar productos (nuevo GET endpoint - REST correcto)
curl -X GET "http://localhost:8080/products/search?name=iPhone&page=0&size=10"

# Buscar con múltiples filtros
curl -X GET "http://localhost:8080/products/search?name=iPhone&price=999.99&stock=5&page=0&size=20"

# Buscar solo con paginación (listar todos)
curl -X GET "http://localhost:8080/products/search?page=0&size=50"

# Crear producto
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

**¡Documentación completa y profesional lista para usar! 🚀**

> Para cualquier duda o mejora, contacta a **Wembie** - El arquitecto de esta API 😎