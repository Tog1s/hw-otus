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
		Level: "DEBUG",
	}
	logger := New(cfg, buffer)

	t.Run("test info", func(t *testing.T) {
		logger.Info(testString)
		require.Contains(t, buffer.String(), testString)
		require.Contains(t, buffer.String(), "INFO")
	})

	t.Run("test error", func(t *testing.T) {
		buffer.Reset()
		logger.Error(testString)
		require.Contains(t, buffer.String(), testString)
		require.Contains(t, buffer.String(), "ERROR")
	})

	t.Run("test warn", func(t *testing.T) {
		buffer.Reset()
		logger.Warn(testString)
		require.Contains(t, buffer.String(), testString)
		require.Contains(t, buffer.String(), "WARN")
	})

	t.Run("test debug", func(t *testing.T) {
		buffer.Reset()
		logger.Debug(testString)
		require.Contains(t, buffer.String(), testString)
		require.Contains(t, buffer.String(), "DEBUG")
	})
}
