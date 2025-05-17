package errors

import (
	"log/slog"
	"os"
)

func SetupDefaultLogger(opts *slog.HandlerOptions) *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	return logger
}

func SetLoggerLevel(level string) *slog.HandlerOptions {
	opts := &slog.HandlerOptions{}

	switch level {
	case "INFO":
		opts.AddSource = false
		opts.Level = slog.LevelInfo
	case "DEBUG":
		opts.AddSource = true
		opts.Level = slog.LevelDebug
	}

	return opts
}
