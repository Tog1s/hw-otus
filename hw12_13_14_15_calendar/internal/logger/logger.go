package logger

import (
	"io"
	"log/slog"

	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/config"
)

type AppLogger struct {
	logger *slog.Logger
}

func New(cfg config.Logger, output io.Writer) *AppLogger {
	logLevel := parseLogLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewJSONHandler(output, opts)
	logger := slog.New(handler)

	return &AppLogger{
		logger: logger,
	}
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "INFO":
		return slog.LevelInfo
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	case "DEBUG":
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func (l *AppLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *AppLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *AppLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *AppLogger) Debug(msg string) {
	l.logger.Debug(msg)
}
