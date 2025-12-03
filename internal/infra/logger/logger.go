package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// InitLogger initializes the global logger
func InitLogger(level, format string) {
	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// Debug logs a debug message with structured fields
func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

// Info logs an info message with structured fields
func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

// Warn logs a warning message with structured fields
func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

// Error logs an error message with structured fields
func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}
