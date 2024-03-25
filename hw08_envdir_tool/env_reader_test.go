package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Validate environment variable", func(t *testing.T) {
		envs := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}
		envFromDir, err := ReadDir("testdata/env")

		require.Equal(t, err, nil)
		require.Equal(t, envs, envFromDir)
	})
}
