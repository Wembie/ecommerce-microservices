package config

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ExtractTraceIDFromGin(c *gin.Context) string {
	priorityHeaders := []string{
		"X-Trace-Id",
		"X-Request-Id",
		"X-Correlation-Id",
		"Trace-Id",
		"Request-Id",
	}

	for _, header := range priorityHeaders {
		if value := c.GetHeader(header); value != "" {
			return value
		}
	}

	return GenerateTraceID()
}

func GenerateTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback-" + strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return fmt.Sprintf("%x", b)
}

func CreateLoggerWithTraceID(baseLogger *zap.Logger, ctx context.Context) *zap.Logger {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		traceID := ExtractTraceIDFromGin(ginCtx)
		return baseLogger.With(zap.String("trace_id", traceID))
	}
	
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return baseLogger.With(zap.String("trace_id", traceID))
	}
	
	return baseLogger.With(zap.String("trace_id", GenerateTraceID()))
}