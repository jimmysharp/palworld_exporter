package log

import (
	"log/slog"
	"os"

	"github.com/jimmysharp/palworld_exporter/config"
)

func NewLogger(config *config.Config) *slog.Logger {
	var logLevel slog.Level
	switch config.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		panic("invalid log level")
	}

	var handler slog.Handler
	switch config.LogFormat {
	case "default":
		handler = slog.Default().Handler()
		slog.SetLogLoggerLevel(logLevel)
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	default:
		panic("invalid log format")
	}

	logger := slog.New(handler)
	return logger
}
