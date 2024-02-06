package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdEnv := make([]string, 0, len(cmd))
	for name, value := range env {
		if value.NeedRemove {
			os.Unsetenv(name)
			continue
		}
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", name, value.Value))
	}

	cmdRun := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	cmdRun.Env = append(os.Environ(), cmdEnv...)
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	cmdRun.Stdin = os.Stdin

	if err := cmdRun.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
