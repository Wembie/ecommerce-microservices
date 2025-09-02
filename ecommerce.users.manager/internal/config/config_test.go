package config_test

import (
	"ecommerce.users.manager/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Run("GetConfig with all valid environment variables", func(t *testing.T) {
		envVars := map[string]string{
			"DB_HOST":       "localhost",
			"DB_PORT":       "5432",
			"DB_USER":       "testuser",
			"DB_PASSWORD":   "testpass",
			"DB_NAME":       "testdb",
			"DB_INSECURE":   "true",
			"MAX_IDLE_CONN": "20",
			"MAX_OPEN_CONN": "15",
			"APP_NAME":      "test-app",
			"APP_PORT":      "50051",
			"JWT_SECRET":    "testsecret",
		}

		setEnvVars(t, envVars)
		defer unsetEnvVars(envVars)

		conf := config.GetConfig()

		assert.Equal(t, "localhost", conf.DBHost)
		assert.Equal(t, "5432", conf.DBPort)
		assert.Equal(t, "testuser", conf.DBUser)
		assert.Equal(t, "testpass", conf.DBPassword)
		assert.Equal(t, "testdb", conf.DBName)
		assert.True(t, conf.DBInsecure)
		assert.Equal(t, 20, conf.MaxConnection)
		assert.Equal(t, 15, conf.MaxOpenConnection)
		assert.Equal(t, "test-app", conf.ApplicationName)
		assert.Equal(t, "50051", conf.AppPort)
	})

	t.Run("GetConfig with default values", func(t *testing.T) {
		envVars := map[string]string{
			"DB_HOST":     "localhost",
			"DB_PORT":     "5432",
			"DB_USER":     "testuser",
			"DB_PASSWORD": "testpass",
			"DB_NAME":     "testdb",
			"APP_NAME":    "test-app",
			"APP_PORT":    "50051",
			"JWT_SECRET":    "testsecret",
		}

		setEnvVars(t, envVars)
		defer unsetEnvVars(envVars)

		conf := config.GetConfig()

		assert.Equal(t, 10, conf.MaxConnection)
		assert.Equal(t, 10, conf.MaxOpenConnection)
		assert.False(t, conf.DBInsecure)
		assert.Equal(t, "50051", conf.AppPort)
	})

	t.Run("GetConfig should panic on missing required DB_HOST", func(t *testing.T) {
		envVars := map[string]string{
			"DB_PORT":     "5432",
			"DB_USER":     "testuser",
			"DB_PASSWORD": "testpass",
			"DB_NAME":     "testdb",
			"APP_NAME":    "test-app",
			"APP_PORT":    "50051",
		}

		setEnvVars(t, envVars)
		defer unsetEnvVars(envVars)

		assert.Panics(t, func() {
			config.GetConfig()
		})
	})

	t.Run("GetConfig should panic on invalid MAX_IDLE_CONN", func(t *testing.T) {
		envVars := map[string]string{
			"DB_HOST":       "localhost",
			"DB_PORT":       "5432",
			"DB_USER":       "testuser",
			"DB_PASSWORD":   "testpass",
			"DB_NAME":       "testdb",
			"APP_NAME":      "test-app",
			"APP_PORT":      "50051",
			"MAX_IDLE_CONN": "invalid",
		}

		setEnvVars(t, envVars)
		defer unsetEnvVars(envVars)

		assert.Panics(t, func() {
			config.GetConfig()
		})
	})

	t.Run("GetConfig should panic on invalid MAX_OPEN_CONN", func(t *testing.T) {
		envVars := map[string]string{
			"DB_HOST":       "localhost",
			"DB_PORT":       "5432",
			"DB_USER":       "testuser",
			"DB_PASSWORD":   "testpass",
			"DB_NAME":       "testdb",
			"APP_NAME":      "test-app",
			"APP_PORT":      "50051",
			"MAX_OPEN_CONN": "invalid",
		}

		setEnvVars(t, envVars)
		defer unsetEnvVars(envVars)

		assert.Panics(t, func() {
			config.GetConfig()
		})
	})
}

func setEnvVars(t *testing.T, envVars map[string]string) {
	for key, value := range envVars {
		err := os.Setenv(key, value)
		assert.NoError(t, err)
	}
}

func unsetEnvVars(envVars map[string]string) {
	for key := range envVars {
		os.Unsetenv(key)
	}
}
