package app

import (
	"log"
	"log/slog"
	"os"
)

func mustSetupLogger(level string) *slog.Logger {
	var logger *slog.Logger

	switch level {
	case "debug":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case "info":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	case "warn":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	case "error":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	default:
		log.Fatalf("invalid logger level: %s", level)
	}

	return logger
}
