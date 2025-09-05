package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort            string
	UserManagerHost    string
}

func GetConfig() *Config {
	return &Config{
		AppPort:         getEnvOrPanic("APP_PORT"),
		UserManagerHost: getEnvOrPanic("USER_MANAGER_HOST"),
	}
}

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || len(value) == 0 {
		panic(fmt.Sprintf("Environment variable %s not found", key))
	}
	return value
}