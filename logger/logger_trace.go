package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	TraceIDKey contextKey = "trace_id"
	userIDKey  contextKey = "user_id"
)

// context key constants
const (
	ContextKeyLogger = "logger" // logger context key
)

// WithTraceID adds traceID to context and returns a new context with logger
func WithTraceID(ctx context.Context, traceID string) context.Context {
	// Add traceID to context
	ctx = context.WithValue(ctx, TraceIDKey, traceID)

	// Get existing logger or create new one
	logger := FromContext(ctx)

	// Add traceID to logger
	logger = &Logger{
		SugaredLogger: logger.With(zap.String(string(TraceIDKey), traceID)),
	}

	// Add logger to context
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

// WithUserID adds userID to context and returns a new context with logger
func WithUserID(ctx context.Context, userID string) context.Context {
	// Add userID to context
	ctx = context.WithValue(ctx, userIDKey, userID)

	// Get existing logger or create new one
	logger := FromContext(ctx)

	// Add userID to logger
	logger = &Logger{
		SugaredLogger: logger.With(zap.String(string(userIDKey), userID)),
	}

	// Add logger to context
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

// WithContext adds both traceID and userID to context
func WithContext(ctx context.Context, traceID string, userID string) context.Context {
	// Add traceID first
	ctx = WithTraceID(ctx, traceID)
	// Then add userID
	return WithUserID(ctx, userID)
}

// FromContext retrieves logger from context
func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ContextKeyLogger).(*Logger); ok {
		return l
	}

	return DefaultLogger
}

// GetTraceID gets traceID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// GetUserID gets userID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

func GetRequestID(ctx context.Context) string {
	return GetTraceID(ctx)
}
