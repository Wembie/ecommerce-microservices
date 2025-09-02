DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'ecommerce_db') THEN
      CREATE DATABASE ecommerce_db;
   END IF;
END
$$;

DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'app_user') THEN
      CREATE USER app_user WITH PASSWORD 'app_password';
   END IF;
END
$$;

GRANT ALL PRIVILEGES ON DATABASE ecommerce_db TO app_user;