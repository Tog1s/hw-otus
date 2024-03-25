package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Status Code is zero", func(t *testing.T) {
		cmd := []string{"testdata/echo.sh", "arg1", "arg2"}
		envs := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		statusCode := RunCmd(cmd, envs)
		require.Equal(t, 0, statusCode)
	})

	t.Run("Arguments error", func(t *testing.T) {
		cmd := []string{"testdata/echo.sj"}
		statusCode := RunCmd(cmd, nil)

		require.Equal(t, 1, statusCode)
	})
}
