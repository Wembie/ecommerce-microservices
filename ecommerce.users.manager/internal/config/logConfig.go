package config

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func extractTraceID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
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
	traceID := extractTraceID(ctx)
	return baseLogger.With(zap.String("trace_id", traceID))
}
