-- Crear bases de datos para cada microservicio
CREATE DATABASE user_service;
CREATE DATABASE product_service;
CREATE DATABASE order_service;

-- Opcional: crear un usuario común
CREATE USER app_user WITH PASSWORD 'app_password';
GRANT ALL PRIVILEGES ON DATABASE user_service TO app_user;
GRANT ALL PRIVILEGES ON DATABASE product_service TO app_user;
GRANT ALL PRIVILEGES ON DATABASE order_service TO app_user;