// logging provides structured helpers for request/error logs
package logging

import (
	"log/slog"
	"os"
)

// LevelFromString converts .env LOG_LEVEL string to slog.Level
func LevelFromString(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewLogger returns a JSON slog logger with the given level
func NewLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

// RequestStart logs method, path, request_id at the start of a request
func RequestStart(log *slog.Logger, method, path, requestID string) {
	log.Info("request start",
		"method", method,
		"path", path,
		"request_id", requestID,
	)
}

// RequestEnd logs status and duration when the request finishes
func RequestEnd(log *slog.Logger, method, path, requestID string, status int, durationMs int64) {
	log.Info("request end",
		"method", method,
		"path", path,
		"request_id", requestID,
		"status", status,
		"duration_ms", durationMs,
	)
}

// LogError writes an error with optional key-value arguments
func LogError(log *slog.Logger, err error, msg string, args ...any) {
	a := append([]any{"error", err}, args...)
	log.Error(msg, a...)
}
