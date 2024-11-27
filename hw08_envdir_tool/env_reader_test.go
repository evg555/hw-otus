package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	envDir = "testdata/env"
)

func TestReadDir(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		expected := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: "\"hello\""},
			"UNSET": EnvValue{NeedRemove: true},
		}

		actual, err := ReadDir(envDir)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("dir path not exist", func(t *testing.T) {
		actual, err := ReadDir("")
		require.Error(t, os.ErrNotExist, err)
		require.Nil(t, actual)
	})

	t.Run("empty env dir", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(envDir, "test")
		if err != nil {
			t.Error(err)
			return
		}
		defer os.RemoveAll(tempDir)

		expected := Environment{}

		actual, err := ReadDir(tempDir)
		require.Error(t, os.ErrNotExist, err)
		require.Equal(t, expected, actual)
	})

	t.Run("dir inside env dir", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(envDir, "test")
		if err != nil {
			t.Error(err)
			return
		}
		defer os.RemoveAll(tempDir)

		expected := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: "\"hello\""},
			"UNSET": EnvValue{NeedRemove: true},
		}

		actual, err := ReadDir(envDir)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
