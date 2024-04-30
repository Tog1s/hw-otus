package logger

import (
	"fmt"
	"io"

	"github.com/tog1s/hw-otus/hw12_13_14_15_calendar/internal/config"
)

type Logger struct {
}

func New(cfg config.Logger, output io.Writer) *Logger {

	return &Logger{}
}

func (l Logger) Info(msg string) {
	fmt.Printf("[INFO] %s", msg)
}

func (l Logger) Error(msg string) {
	fmt.Println(msg)
}

func (l Logger) Warn(msg string) {
	fmt.Println(msg)
}

func (l Logger) Debug(msg string) {
	fmt.Println(msg)
}
