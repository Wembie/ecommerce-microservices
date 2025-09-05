package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBInsecure          bool
	MaxConnection       int
	MaxOpenConnection   int
	ApplicationName     string
	AppPort             string
	UserManagerHost     string
	ProductManagerHost  string
}

func GetConfig() *Config {
	maxConnection, err := strconv.Atoi(getEnvOrDefault("MAX_IDLE_CONN", "10"))
	if err != nil {
		panic(fmt.Sprintf("Invalid MAX_IDLE_CONN value: %s", err))
	}
	maxOpenConnection, err := strconv.Atoi(getEnvOrDefault("MAX_OPEN_CONN", "10"))
	if err != nil {
		panic(fmt.Sprintf("Invalid MAX_OPEN_CONN value: %s", err))
	}

	return &Config{
		DBHost:              getEnvOrPanic("DB_HOST"),
		DBPort:              getEnvOrPanic("DB_PORT"),
		DBUser:              getEnvOrPanic("DB_USER"),
		DBPassword:          getEnvOrPanic("DB_PASSWORD"),
		DBName:              getEnvOrPanic("DB_NAME"),
		DBInsecure:          getEnvOrDefault("DB_INSECURE", "false") == "true",
		MaxConnection:       maxConnection,
		MaxOpenConnection:   maxOpenConnection,
		ApplicationName:     getEnvOrPanic("APP_NAME"),
		AppPort:             getEnvOrPanic("APP_PORT"),
		UserManagerHost:     getEnvOrDefault("USER_MANAGER_HOST", "dns:localhost:50051"),
		ProductManagerHost:  getEnvOrDefault("PRODUCT_MANAGER_HOST", "http://localhost:8081"),
	}
}

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || len(value) == 0 {
		panic(fmt.Sprintf("Environment variable %s not found", key))
	}
	return value
}

func getEnvOrDefault(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok || len(value) == 0 {
		return def
	}
	return value
}