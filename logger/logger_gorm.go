package logger

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger.Interface
	SlowThreshold time.Duration
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: time.Second,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	log := FromContext(ctx)
	log.Infof(msg, data...)
}

// Warn print warning
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log := FromContext(ctx)
	log.Warnf(msg, data...)
}

// Error print error
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	log := FromContext(ctx)
	log.Errorf(msg, data...)
}

// Trace print sql and log level
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// get traceID and userID
	traceID := GetTraceID(ctx)
	userID := GetUserID(ctx)

	// get log instance
	log := FromContext(ctx)

	// build log fields
	fields := []interface{}{
		"elapsed", elapsed,
		"rows", rows,
		"sql", sql,
	}

	if traceID != "" {
		fields = append(fields, "trace_id", traceID)
	}

	if userID != "" {
		fields = append(fields, "user_id", userID)
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		fields = append(fields, "error", err)
		log.Errorw("SQL Error", fields...)
	case elapsed > l.SlowThreshold:
		log.Warnw("Slow SQL", fields...)
	default:
		log.Debugw("SQL Trace", fields...)
	}
}
