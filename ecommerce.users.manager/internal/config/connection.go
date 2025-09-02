package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func GetConnection(config *Config) *bun.DB {
	addr := fmt.Sprintf("%s:%s", config.DBHost, config.DBPort)

	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(addr),
		pgdriver.WithUser(config.DBUser),
		pgdriver.WithPassword(config.DBPassword),
		pgdriver.WithDatabase(config.DBName),
		pgdriver.WithApplicationName(config.ApplicationName),
		pgdriver.WithInsecure(config.DBInsecure),
	)

	sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUN_LOG_LEVEL"),
	))

	db.SetMaxOpenConns(config.MaxOpenConnection)
	db.SetMaxIdleConns(config.MaxConnection)

	db.SetConnMaxIdleTime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(fmt.Sprintf("cannot connect to database: %s", err))
	}

	return db
}
