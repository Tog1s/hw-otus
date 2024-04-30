package logger

import (
	"bytes"
	"testing"

	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/config"
)

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	cfg := config.Logger{
		Level: "INFO",
	}
	t.Run("test logger", func(t *testing.T) {
		log := New(cfg, buffer)
		log.Info("info message")
	})
}
