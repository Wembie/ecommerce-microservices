@echo off
chcp 65001 >nul
REM =================================================================
REM Script to generate Swagger documentation
REM Created by: Wembie
REM =================================================================

REM Change to the project root directory
cd /d "%~dp0.."

echo.
echo ===============================================
echo  Generating Swagger documentation for Products Manager
echo ===============================================
echo.

REM Check if swag is installed
swag --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] 'swag' is not installed
    echo.
    echo [INFO] To install it, run:
    echo    go install github.com/swaggo/swag/cmd/swag@latest
    echo.
    pause
    exit /b 1
)

echo [OK] swag is installed correctly
echo.

REM Clean previous documentation
echo [INFO] Cleaning up previous documentation...
if exist "docs\docs.go" del /q "docs\docs.go"
if exist "docs\swagger.json" del /q "docs\swagger.json"
if exist "docs\swagger.yaml" del /q "docs\swagger.yaml"
echo.

REM Generate documentation
echo [INFO] Generating new documentation...
echo.
swag init -g cmd/main.go -o docs --parseDependency --parseInternal

if %errorlevel% equ 0 (
    echo.
    echo [SUCCESS] Documentation generated successfully!
    echo.
    echo [INFO] Generated files:
    echo    - docs/docs.go      (Go file with Swagger data)
    echo    - docs/swagger.json (OpenAPI specification in JSON)
    echo    - docs/swagger.yaml (OpenAPI specification in YAML)
    echo.
    echo [INFO] To view the documentation:
    echo    1. Run the application: go run cmd/main.go
    echo    2. Visit: http://localhost:8080/swagger/index.html
    echo.
    echo [TIP] Documentation updates automatically with each build
    echo.
) else (
    echo.
    echo [ERROR] Failed to generate documentation
    echo.
)

pause
