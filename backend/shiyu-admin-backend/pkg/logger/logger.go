package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey int

const (
	traceIDKey ctxKey = iota
)

var defaultLogger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

// Init initializes the global logger with given level/format and optional file retention.
// If retentionDays > 0, logs are also written to a rolling file with that many days retained.
// maxSizeMB controls the maximum size of a single log file in megabytes (if <=0, a default is used).
func Init(level, format string, retentionDays int, filePath string, maxSizeMB int) {
	lvl := parseLevel(level)
	opts := &slog.HandlerOptions{Level: lvl}

	var writer io.Writer = os.Stdout
	if retentionDays > 0 {
		if filePath == "" {
			filePath = "logs/app.log"
		}
		if dir := filepath.Dir(filePath); dir != "" && dir != "." {
			_ = os.MkdirAll(dir, 0o755)
		}
		if maxSizeMB <= 0 {
			maxSizeMB = 100
		}
		lj := &lumberjack.Logger{
			Filename:   filePath,
			MaxAge:     retentionDays,
			MaxSize:    maxSizeMB,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   false,
		}
		writer = io.MultiWriter(os.Stdout, lj)
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "", "json":
		handler = slog.NewJSONHandler(writer, opts)
	case "text", "console":
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewJSONHandler(writer, opts)
	}

	defaultLogger = slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info", "":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithTraceID stores the trace ID into the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, traceIDKey, traceID)
}

func traceIDFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	if v := ctx.Value(traceIDKey); v != nil {
		if id, ok := v.(string); ok && id != "" {
			return id, true
		}
	}
	return "", false
}

// L returns the underlying global logger.
func L() *slog.Logger {
	return defaultLogger
}

// Debug logs a debug level message with optional key-value pairs.
func Debug(ctx context.Context, msg string, args ...any) {
	logWithLevel(ctx, slog.LevelDebug, msg, args...)
}

// Info logs an info level message with optional key-value pairs.
func Info(ctx context.Context, msg string, args ...any) {
	logWithLevel(ctx, slog.LevelInfo, msg, args...)
}

// Warn logs a warning level message with optional key-value pairs.
func Warn(ctx context.Context, msg string, args ...any) {
	logWithLevel(ctx, slog.LevelWarn, msg, args...)
}

// Error logs an error level message with optional key-value pairs.
func Error(ctx context.Context, msg string, args ...any) {
	logWithLevel(ctx, slog.LevelError, msg, args...)
}

func logWithLevel(ctx context.Context, level slog.Level, msg string, args ...any) {
	if defaultLogger == nil {
		return
	}

	if traceID, ok := traceIDFromContext(ctx); ok {
		args = append([]any{"trace_id", traceID}, args...)
	}

	switch level {
	case slog.LevelDebug:
		defaultLogger.Debug(msg, args...)
	case slog.LevelInfo:
		defaultLogger.Info(msg, args...)
	case slog.LevelWarn:
		defaultLogger.Warn(msg, args...)
	case slog.LevelError:
		defaultLogger.Error(msg, args...)
	default:
		defaultLogger.Info(msg, args...)
	}
}
