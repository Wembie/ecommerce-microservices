package config_test

import (
	"os"
	"testing"

	"ecommerce.auth/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Run("GetConfig with all valid environment variables", func(t *testing.T) {
		envVars := map[string]string{
			"APP_PORT":         "8080",
			"USER_MANAGER_HOST": "dns:localhost:50051",
		}

		setEnvVars(t, envVars)
		defer unsetEnvVars(envVars)

		conf := config.GetConfig()

		assert.Equal(t, "8080", conf.AppPort)
		assert.Equal(t, "dns:localhost:50051", conf.UserManagerHost)
	})

	t.Run("GetConfig should panic on missing APP_PORT", func(t *testing.T) {
		envVars := map[string]string{
			"USER_MANAGER_HOST": "dns:localhost:50051",
		}

		setEnvVars(t, envVars)
		defer unsetEnvVars(envVars)

		assert.Panics(t, func() {
			config.GetConfig()
		})
	})

	t.Run("GetConfig should panic on missing USER_MANAGER_HOST", func(t *testing.T) {
		envVars := map[string]string{
			"APP_PORT": "8080",
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
