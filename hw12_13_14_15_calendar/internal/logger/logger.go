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

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(output, opts)
	logger := slog.New(handler)

	return &AppLogger{
		logger: logger,
	}
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
