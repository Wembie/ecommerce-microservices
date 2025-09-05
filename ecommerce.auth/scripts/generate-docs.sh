#!/bin/bash
# =================================================================
# Script to generate Swagger documentation
# Created by: Wembie
# =================================================================

# Change to the project root directory
cd "$(dirname "$0")/.."

echo ""
echo "==============================================="
echo "  Generating Swagger documentation for Products Manager"
echo "==============================================="
echo ""

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "[ERROR] 'swag' is not installed"
    echo ""
    echo "[INFO] To install it, run:"
    echo "   go install github.com/swaggo/swag/cmd/swag@latest"
    echo ""
    exit 1
fi

echo "[OK] swag is installed correctly"
echo ""

# Clean previous documentation
echo "[INFO] Cleaning up previous documentation..."
rm -f docs/docs.go docs/swagger.json docs/swagger.yaml
echo ""

# Generate documentation
echo "[INFO] Generating new documentation..."
echo ""
swag init -g cmd/main.go -o docs --parseDependency --parseInternal

if [ $? -eq 0 ]; then
    echo ""
    echo "[SUCCESS] Documentation generated successfully!"
    echo ""
    echo "[INFO] Generated files:"
    echo "   - docs/docs.go      (Go file with Swagger data)"
    echo "   - docs/swagger.json (OpenAPI specification in JSON)"
    echo "   - docs/swagger.yaml (OpenAPI specification in YAML)"
    echo ""
    echo "[INFO] To view the documentation:"
    echo "   1. Run the application: go run cmd/main.go"
    echo "   2. Visit: http://localhost:8081/swagger/index.html"
    echo ""
    echo "[TIP] The documentation updates automatically with each build"
    echo ""
else
    echo ""
    echo "[ERROR] Failed to generate documentation"
    echo ""
    exit 1
fi
