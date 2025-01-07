package app

import (
	"context"
	"time"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, event storage.Event) error
	ListEventsForDay(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListEventsForWeek(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListEventsForMonth(ctx context.Context, date time.Time) ([]*storage.Event, error)
	Close(ctx context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id string, title string) error {
	event := storage.Event{
		ID:    id,
		Title: title,
	}

	return a.storage.CreateEvent(ctx, event)
}

// TODO
