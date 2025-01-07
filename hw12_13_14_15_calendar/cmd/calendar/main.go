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
	internalhttp "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	internalstorage "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

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

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(cfg, logg, calendar)

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

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
