package main

import (
	"context"
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
	internalstorage "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
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

	var storage app.Storage

	switch cfg.App.Storage {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		storage = sqlstorage.New(cfg.Database)
	default:
		panic(fmt.Sprintf("%s: %v: %s", cfg.App.Storage, internalstorage.ErrStorageNotExist, cfg.App.Storage))
	}

	queue := rabbit.New(cfg.Rabbit)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := storage.Close(ctx); err != nil {
			logg.Error("failed to close connection to storage: " + err.Error())
		}

		if err := queue.Close(); err != nil {
			logg.Error("failed to close connection to queue: " + err.Error())
		}

		os.Exit(0)
	}()

	logg.Info("scheduler is running...")

	for {
		logg.Debug("handling events...")
		now := time.Now()

		err := storage.DeleteOldEvents(ctx, now)
		if err != nil {
			logg.Error("failed to delete old events: " + err.Error())
			return
		}

		events, err := storage.ListEventsForNotify(ctx, now)
		if err != nil {
			logg.Error("failed to list events for notify: " + err.Error())
			return
		}

		for _, event := range events {
			var userID string

			if event.UserID.Valid {
				userID = event.UserID.String
			}

			notification := app.Notification{
				EventID: event.ID,
				Title:   event.Title,
				Date:    event.StartDate.Format("2006-01-02 15:04"),
				UserID:  userID,
			}

			err = queue.Add(notification)
			if err != nil {
				logg.Error("failed to add event to queue: " + err.Error())
			}
		}

		logg.Debug("sleeping 1 day...")
		time.Sleep(24 * time.Hour)
	}
}
