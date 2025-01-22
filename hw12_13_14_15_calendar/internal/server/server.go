package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/config"
	internalgrpc "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/server/http"
)

const (
	httpServer = "http"
	grpcServer = "grpc"
)

var ErrServerNotExist = errors.New("server not exist")

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, id string, title string) error
	UpdateEvent(ctx context.Context, id string, event app.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context, date, period string) ([]app.Event, error)
}

func New(cfg config.Config, logger Logger, app Application) Server {
	var server Server

	switch cfg.App.Server {
	case httpServer:
		server = internalhttp.NewServer(cfg, logger, app)
	case grpcServer:
		server = internalgrpc.NewServer(cfg, logger, app)
	default:
		panic(fmt.Sprintf("%s: %v: %s", cfg.App.Storage, ErrServerNotExist, cfg.App.Server))
	}

	return server
}
