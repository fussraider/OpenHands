package logger

import (
	"log/slog"
	"os"
	"strings"
)

func Init() {
	level := slog.LevelInfo
	envLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if envLevel == "debug" || os.Getenv("DEBUG") == "true" {
		level = slog.LevelDebug
	} else if envLevel == "warn" || envLevel == "warning" {
		level = slog.LevelWarn
	} else if envLevel == "error" {
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
