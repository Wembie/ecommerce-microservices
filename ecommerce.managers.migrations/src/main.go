package main

import (
	"embed"

	"ecommerce.managers.migrations/src/config"
	"ecommerce.managers.migrations/src/db"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	config.Logger.Info("🚀 Starting migration")

	cfg := config.InitEnv()
	
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		config.Logger.Fatal("❌ Required database environment variables are missing")
	}

	instance := db.GetConnection()
	defer func() {
		if err := instance.DB.Close(); err != nil {
			config.Logger.Error("❌ Error closing connection", zap.Error(err))
		}
	}()

	db.Migrate(instance.DB, embedMigrations)

	config.Logger.Info("✅ Migration finished successfully")
}