package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("command not exist", func(t *testing.T) {
		tempFile, err := os.CreateTemp(envDir, "stdout")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		cmd := []string{"not_exist", "arg1", "arg2"}
		env := Environment{}
		outExpected := "exec: \"not_exist\": executable file not found in $PATH"

		stdout := os.Stdout
		os.Stdout = tempFile

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 1, exitCode)

		os.Stdout = stdout

		_, err = tempFile.Seek(0, 0)
		require.NoError(t, err)
		out, err := io.ReadAll(tempFile)
		require.NoError(t, err)

		require.Contains(t, string(out), outExpected)
	})

	t.Run("command returns error", func(t *testing.T) {
		tempFile, err := os.CreateTemp(envDir, "stdout")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		cmd := []string{"rm", "/root"}
		env := Environment{}
		outExpected := "rm: cannot remove '/root': Is a directory"

		stderr := os.Stderr
		os.Stderr = tempFile

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 1, exitCode)

		os.Stderr = stderr

		_, err = tempFile.Seek(0, 0)
		require.NoError(t, err)
		out, err := io.ReadAll(tempFile)
		require.NoError(t, err)

		require.Contains(t, string(out), outExpected)
	})

	t.Run("empty command", func(t *testing.T) {
		var cmd []string
		env := Environment{}

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)
	})

	t.Run("empty command params", func(t *testing.T) {
		tempFile, err := os.CreateTemp(envDir, "stdout")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		cmd := []string{"echo"}
		env := Environment{
			"FOO": EnvValue{
				Value: "bar",
			},
		}
		outExpected := "\n"

		stdout := os.Stdout
		os.Stdout = tempFile

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)

		os.Stdout = stdout

		_, err = tempFile.Seek(0, 0)
		require.NoError(t, err)
		out, err := io.ReadAll(tempFile)
		require.NoError(t, err)

		require.Equal(t, outExpected, string(out))
	})

	t.Run("with replace and remove env values", func(t *testing.T) {
		tempFile, err := os.CreateTemp(envDir, "stdout")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		os.Setenv("HELLO", "SHOULD_REPLACE")
		os.Setenv("FOO", "SHOULD_REPLACE")
		os.Setenv("UNSET", "SHOULD_REMOVE")
		os.Setenv("ADDED", "from original env")
		os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

		cmd := []string{"./testdata/echo.sh", "arg1=1", "arg2=2"}
		env := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: "\"hello\""},
			"UNSET": EnvValue{NeedRemove: true},
		}
		outExpected := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`

		stdout := os.Stdout
		os.Stdout = tempFile

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)

		os.Stdout = stdout

		_, err = tempFile.Seek(0, 0)
		require.NoError(t, err)
		out, err := io.ReadAll(tempFile)
		require.NoError(t, err)

		require.Contains(t, string(out), outExpected)
	})
}
