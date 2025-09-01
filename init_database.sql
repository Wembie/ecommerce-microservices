-- Crear bases de datos para cada microservicio
CREATE DATABASE ecommerce_db;

-- Opcional: crear un usuario común
CREATE USER app_user WITH PASSWORD 'app_password';
GRANT ALL PRIVILEGES ON DATABASE ecommerce_db TO app_user;