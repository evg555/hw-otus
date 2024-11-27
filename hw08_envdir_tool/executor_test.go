package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("command not exist", func(t *testing.T) {
		cmd := []string{"not_exist", "arg1", "arg2"}
		env := Environment{}
		outExpected := "exec: \"not_exist\": executable file not found in $PATH\n"

		var buf bytes.Buffer
		originalStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 1, exitCode)

		w.Close()
		os.Stdout = originalStdout

		_, err := buf.ReadFrom(r)
		if err != nil {
			t.Error()
		}
		output := buf.String()
		require.Equal(t, outExpected, output)
	})

	t.Run("command returns error", func(t *testing.T) {
		cmd := []string{"rm", "/root"}
		env := Environment{}
		outExpected := "rm: cannot remove '/root': Is a directory\n\n"

		var buf bytes.Buffer
		originalStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 1, exitCode)

		w.Close()
		os.Stdout = originalStdout

		_, err := buf.ReadFrom(r)
		if err != nil {
			t.Error()
		}
		output := buf.String()
		require.Equal(t, outExpected, output)
	})

	t.Run("empty command", func(t *testing.T) {
		var cmd []string
		env := Environment{}

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)
	})

	t.Run("empty command params", func(t *testing.T) {
		cmd := []string{"echo"}
		env := Environment{
			"FOO": EnvValue{
				Value: "bar",
			},
		}
		outExpected := "\n\n"

		var buf bytes.Buffer
		originalStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)

		w.Close()
		os.Stdout = originalStdout

		_, err := buf.ReadFrom(r)
		if err != nil {
			t.Error()
		}
		output := buf.String()
		require.Equal(t, outExpected, output)
	})

	t.Run("with replace and remove env values", func(t *testing.T) {
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

		var buf bytes.Buffer
		originalStdout := os.Stdout

		r, w, _ := os.Pipe()
		os.Stdout = w

		exitCode := RunCmd(cmd, env)
		require.Equal(t, 0, exitCode)

		w.Close()
		os.Stdout = originalStdout

		_, err := buf.ReadFrom(r)
		if err != nil {
			t.Error()
		}
		output := buf.String()
		require.Equal(t, outExpected, output)
	})
}
