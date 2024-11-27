package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return
	}

	var args []string

	if len(cmd) > 1 {
		args = cmd[1:]
	}

	command := exec.Command(cmd[0], args...) //nolint:gosec

	for key, envValue := range env {
		err := os.Unsetenv(key)
		if err != nil {
			fmt.Println(err)
			return 1
		}

		if envValue.NeedRemove {
			continue
		}

		command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, envValue.Value))
	}

	command.Env = append(command.Env, os.Environ()...)

	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				fmt.Println(out.String())
				return status.ExitStatus()
			}
		}

		fmt.Println(err)
		return 1
	}

	fmt.Println(out.String())
	return
}
