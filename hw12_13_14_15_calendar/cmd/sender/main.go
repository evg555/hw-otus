package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/queue/rabbit"
)

func main() {
	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	cfg := config.NewConfig()
	logg := logger.New(cfg.Logger)

	queue := rabbit.New(cfg.Rabbit)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := queue.Close(); err != nil {
			logg.Error("failed to close connection to queue: " + err.Error())
		}

		os.Exit(0)
	}()

	logg.Info("sender is running...")

	for rawMessage := range queue.Get() {
		var notification app.Notification

		if err := json.Unmarshal(rawMessage.Body, &notification); err != nil {
			logg.Error("failed to unmarshal raw message: " + err.Error())
			continue
		}

		logg.Info(fmt.Sprintf("sent message for user %s: event '%s' on %s",
			notification.UserID, notification.Title, notification.Date))
	}
}
