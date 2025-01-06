package logger

import (
	"errors"
	"fmt"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/config"
	log "github.com/sirupsen/logrus"
)

var ErrFormatNotExist = errors.New("format not exist")

type Logger struct {
	logger *log.Logger
}

func New(cfg config.LoggerConf) Logger {
	logger := log.New()

	switch cfg.Format {
	case "json":
		logger.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	case "text":
		logger.SetFormatter(&log.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		panic(fmt.Sprintf("init logger error: %v: %s", ErrFormatNotExist, cfg.Format))
	}

	loglevel, err := log.ParseLevel(cfg.Level)
	if err != nil {
		panic(fmt.Sprintf("init logger error: %v", err))
	}

	logger.SetLevel(loglevel)

	return Logger{
		logger: logger,
	}
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
