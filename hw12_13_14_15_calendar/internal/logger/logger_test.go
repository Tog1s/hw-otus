package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/config"
)

var testString string = "test message string"

func TestLogger(t *testing.T) {
	buffer := &bytes.Buffer{}
	cfg := config.Logger{
		Level: "INFO",
	}
	logger := New(cfg, buffer)

	t.Run("test info", func(t *testing.T) {
		logger.Info(testString)
		require.Contains(t, buffer.String(), testString)
	})
}
