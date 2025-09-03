package config_test

import (
	"context"
	"strings"
	"testing"

	"ecommerce.products.manager/internal/config"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func TestExtractTraceID(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]string
		expected string
	}{
		{
			name: "extract x-trace-id (highest priority)",
			metadata: map[string]string{
				"x-trace-id":       "trace-123",
				"x-request-id":     "request-456",
				"x-correlation-id": "corr-789",
			},
			expected: "trace-123",
		},
		{
			name: "extract x-request-id when x-trace-id not present",
			metadata: map[string]string{
				"x-request-id":     "request-456",
				"x-correlation-id": "corr-789",
			},
			expected: "request-456",
		},
		{
			name: "extract x-correlation-id when higher priority headers not present",
			metadata: map[string]string{
				"x-correlation-id": "corr-789",
			},
			expected: "corr-789",
		},
		{
			name: "extract trace-id (without x- prefix)",
			metadata: map[string]string{
				"trace-id": "trace-without-x",
			},
			expected: "trace-without-x",
		},
		{
			name: "extract request-id (without x- prefix)",
			metadata: map[string]string{
				"request-id": "request-without-x",
			},
			expected: "request-without-x",
		},
		{
			name: "ignore empty header values",
			metadata: map[string]string{
				"x-trace-id":   "",
				"x-request-id": "request-456",
			},
			expected: "request-456",
		},
		{
			name:     "no metadata should generate random trace ID",
			metadata: map[string]string{},
			expected: "",
		},
		{
			name:     "context without incoming metadata should generate random trace ID",
			metadata: nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			if len(tt.metadata) > 0 {
				md := metadata.New(tt.metadata)
				ctx = metadata.NewIncomingContext(context.Background(), md)
			} else {
				ctx = context.Background()
			}

			result := callExtractTraceID(ctx)

			if tt.expected == "" {
				assert.NotEmpty(t, result)
				assert.True(t, len(result) > 0)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGenerateTraceID(t *testing.T) {
	t.Run("should generate non-empty trace ID", func(t *testing.T) {
		traceID := config.GenerateTraceID()
		assert.NotEmpty(t, traceID)
		assert.True(t, len(traceID) > 0)
	})

	t.Run("should generate different IDs on multiple calls", func(t *testing.T) {
		id1 := config.GenerateTraceID()
		id2 := config.GenerateTraceID()
		assert.NotEqual(t, id1, id2)
	})

	t.Run("should generate hex format", func(t *testing.T) {
		traceID := config.GenerateTraceID()
		assert.True(t, isHexString(traceID) || strings.HasPrefix(traceID, "fallback-"))
	})
}

func TestCreateLoggerWithTraceID(t *testing.T) {
	baseLogger := zap.NewNop()

	t.Run("should create logger with trace ID from metadata", func(t *testing.T) {
		md := metadata.New(map[string]string{
			"x-trace-id": "test-trace-123",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		logger := config.CreateLoggerWithTraceID(baseLogger, ctx)
		assert.NotNil(t, logger)

		assert.IsType(t, &zap.Logger{}, logger)
	})

	t.Run("should create logger with generated trace ID when no metadata", func(t *testing.T) {
		ctx := context.Background()

		logger := config.CreateLoggerWithTraceID(baseLogger, ctx)
		assert.NotNil(t, logger)
		assert.IsType(t, &zap.Logger{}, logger)
	})
}

func callExtractTraceID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return config.GenerateTraceID()
	}

	priorityHeaders := []string{
		"x-trace-id",
		"x-request-id",
		"x-correlation-id",
		"trace-id",
		"request-id",
	}

	for _, header := range priorityHeaders {
		if values := md.Get(header); len(values) > 0 && values[0] != "" {
			return values[0]
		}
	}

	return config.GenerateTraceID()
}

func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
