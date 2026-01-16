package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	DefaultAppsName = "myapp"
	// Global logger instance
	DefaultLogger *Logger
)

// Logger wraps zap logger with additional functionality
type Logger struct {
	*zap.SugaredLogger
}

// LogConfig defines the configuration for the logger
type LogConfig struct {
	// Log file directory
	FileDir string
	// Maximum size in megabytes of the log file before it gets rotated
	MaxSize int
	// Maximum number of old log files to retain
	MaxBackups int
	// Maximum number of days to retain old log files
	MaxAge int
	// Whether to compress old log files
	Compress bool
	// Minimum logging level
	Level string
	// Whether to also log to console
	EnableConsole bool
	//
	AppName string
	// Console encoder config
	ConsoleEncoderConfig string
}

// DefaultConfig returns the default logging configuration
func DefaultConfig() LogConfig {
	return LogConfig{
		FileDir:              "logs",
		MaxSize:              100,    // 100 MB
		MaxBackups:           30,     // keep 30 old files
		MaxAge:               30,     // 30 days
		Compress:             true,   // compress old files
		Level:                "info", // default log level
		EnableConsole:        true,   // enable console logging in development
		AppName:              DefaultAppsName,
		ConsoleEncoderConfig: "development",
	}
}

// projectRootCallerEncoder shows path relative to project root
func projectRootCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	file := caller.File

	// Find the project root by looking for common markers
	projectMarkers := []string{
		"go.mod",
		"cmd/",
		"internal/",
		"api/",
		"pkg/",
	}

	for _, marker := range projectMarkers {
		if idx := strings.Index(file, marker); idx != -1 {
			// Show path from project root
			relPath := file[idx:]
			relPath = filepath.ToSlash(relPath)
			enc.AppendString(relPath + ":" + strconv.Itoa(caller.Line))
			return
		}
	}

	// Fallback to filename:line
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		enc.AppendString(parts[len(parts)-1] + ":" + strconv.Itoa(caller.Line))
	} else {
		enc.AppendString(caller.String())
	}
}

// InitLogger initializes the global logger with the given configuration
func InitLogger(config LogConfig) error {
	var err error
	// Ensure log directory exists
	if err := os.MkdirAll(config.FileDir, 0755); err != nil {
		fmt.Println("Logger os.MkdirAll error:", err)
		return err
	}

	fmt.Println("Logger config.FileDir", config.FileDir)

	fileRotator := &lumberjack.Logger{
		Filename:   filepath.Join(config.FileDir, config.AppName+".log"),
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// Parse log level
	level := zap.NewAtomicLevel()
	err = level.UnmarshalText([]byte(config.Level))
	if err != nil {
		return err
	}

	// Create cores
	var cores []zapcore.Core

	// File core - use custom formatter
	fileCore := zapcore.NewCore(
		NewFormatterEncoder(),
		zapcore.AddSync(fileRotator),
		level,
	)
	cores = append(cores, fileCore)

	// Console core (if enabled)
	if config.EnableConsole {
		consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
		consoleEncoderConfig.TimeKey = "time"
		consoleEncoderConfig.LevelKey = "level"
		consoleEncoderConfig.NameKey = "logger"
		consoleEncoderConfig.CallerKey = "caller"
		consoleEncoderConfig.MessageKey = "msg"
		consoleEncoderConfig.StacktraceKey = "stacktrace"
		consoleEncoderConfig.LineEnding = zapcore.DefaultLineEnding
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		consoleEncoderConfig.EncodeCaller = projectRootCallerEncoder
		var consoleCore zapcore.Core
		if config.ConsoleEncoderConfig == "development" {
			consoleCore = zapcore.NewCore(
				zapcore.NewConsoleEncoder(consoleEncoderConfig),
				zapcore.AddSync(os.Stdout),
				level,
			)
		} else {
			consoleCore = zapcore.NewCore(
				NewFormatterEncoder(),
				zapcore.AddSync(os.Stdout),
				level,
			)
		}
		cores = append(cores, consoleCore)
	}

	// Create logger with caller skip
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1), // caller skip
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// Create sugared logger
	DefaultLogger = &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}

	return nil
}

// Debug logs a message at debug level
func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

// Debugf logs a formatted message at debug level
func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

// Info logs a message at info level
func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

// Infof logs a formatted message at info level
func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

// Warn logs a message at warn level
func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

// Warnf logs a formatted message at warn level
func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

// Error logs a message at error level
func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

// Errorf logs a formatted message at error level
func Errorf(format string, args ...interface{}) {
	DefaultLogger.Errorf(format, args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

// Fatalf logs a formatted message at fatal level and then calls os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	DefaultLogger.Fatalf(format, args...)
}

// WithFields adds structured fields to the logging context
func WithFields(fields map[string]any) *Logger {
	return &Logger{DefaultLogger.With(fields)}
}

// Sync flushes any buffered log entries
func Sync() error {
	return DefaultLogger.Sync()
}

// Trace logs a message at trace level (mapped to debug since zap doesn't have trace)
func (l *Logger) Trace(v ...interface{}) {
	l.Debug(v...)
}

// Notice logs a message at notice level (mapped to info since zap doesn't have notice)
func (l *Logger) Notice(v ...interface{}) {
	l.Info(v...)
}

// Tracef logs a formatted message at trace level (mapped to debug)
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

// Noticef logs a formatted message at notice level (mapped to info)
func (l *Logger) Noticef(format string, v ...interface{}) {
	l.Infof(format, v...)
}

// CtxTracef logs a formatted message at trace level with context
func (l *Logger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Debugf(format, v...)
}

// CtxDebugf logs a formatted message at debug level with context
func (l *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Debugf(format, v...)
}

// CtxInfof logs a formatted message at info level with context
func (l *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Infof(format, v...)
}

// CtxNoticef logs a formatted message at notice level with context (mapped to info)
func (l *Logger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Infof(format, v...)
}

// CtxWarnf logs a formatted message at warn level with context
func (l *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Warnf(format, v...)
}

// CtxErrorf logs a formatted message at error level with context
func (l *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Errorf(format, v...)
}

// CtxFatalf logs a formatted message at fatal level with context
func (l *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	FromContext(ctx).Fatalf(format, v...)
}
