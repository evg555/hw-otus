package logger

import (
	"bytes"
	"testing"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("level debug", func(t *testing.T) {
		var buf bytes.Buffer

		cfg := config.Config{
			Logger: config.LoggerConf{
				Level:  "debug",
				Format: "text",
			},
		}

		logger := New(cfg.Logger)
		require.NotNil(t, logger)

		logger.logger.Out = &buf
		logger.Debug("this is a debug message")

		output := buf.String()
		require.Contains(t, output, "this is a debug message")
	})

	t.Run("level info", func(t *testing.T) {
		var buf bytes.Buffer

		cfg := config.Config{
			Logger: config.LoggerConf{
				Level:  "info",
				Format: "text",
			},
		}

		logger := New(cfg.Logger)
		require.NotNil(t, logger)

		logger.logger.Out = &buf

		logger.Debug("this should not appear")
		logger.Info("this is an info message")

		output := buf.String()
		require.Contains(t, output, "this is an info message")
	})

	t.Run("level warn", func(t *testing.T) {
		var buf bytes.Buffer

		cfg := config.Config{
			Logger: config.LoggerConf{
				Level:  "warn",
				Format: "text",
			},
		}

		logger := New(cfg.Logger)
		require.NotNil(t, logger)
		logger.logger.Out = &buf

		logger.Debug("this should not appear")
		logger.Info("this should not appear")
		logger.Warn("this is a warning message")

		output := buf.String()
		require.NotContains(t, output, "this should not appear")
		require.Contains(t, output, "this is a warning message")
	})

	t.Run("level error", func(t *testing.T) {
		var buf bytes.Buffer

		cfg := config.Config{
			Logger: config.LoggerConf{
				Level:  "error",
				Format: "text",
			},
		}

		logger := New(cfg.Logger)
		require.NotNil(t, logger)
		logger.logger.Out = &buf

		logger.Debug("this should not appear")
		logger.Info("this should not appear")
		logger.Warn("this should not appear")
		logger.Error("this is an error message")

		output := buf.String()
		require.NotContains(t, output, "this should not appear")
		require.Contains(t, output, "this is an error message")
	})

	t.Run("invalid level", func(t *testing.T) {
		require.Panics(t, func() {
			cfg := config.Config{
				Logger: config.LoggerConf{
					Level:  "invalid",
					Format: "text",
				},
			}

			New(cfg.Logger)
		}, "expected panic for invalid log level, but none occurred")
	})

	t.Run("json format", func(t *testing.T) {
		var buf bytes.Buffer

		cfg := config.Config{
			Logger: config.LoggerConf{
				Level:  "info",
				Format: "json",
			},
		}

		logger := New(cfg.Logger)
		require.NotNil(t, logger)
		logger.logger.Out = &buf

		logger.Info("this is a info message")

		output := buf.String()
		require.Contains(t, output, `"msg":"this is a info message"`)
	})

	t.Run("invalid format", func(t *testing.T) {
		require.Panics(t, func() {
			cfg := config.Config{
				Logger: config.LoggerConf{
					Level:  "info",
					Format: "invalid",
				},
			}

			New(cfg.Logger)
		}, "expected panic for invalid log format, but none occurred")
	})
}
