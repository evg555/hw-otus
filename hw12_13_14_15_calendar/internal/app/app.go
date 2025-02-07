package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

const (
	PeriodDay   = "day"
	PeriodWeek  = "week"
	PeriodMonth = "month"
)

var ErrInvalidPeriod = errors.New("invalid period")

type App struct {
	logger  Logger
	storage Storage
}

//go:generate mockery --name=Logger
type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

//go:generate mockery --name=Storage
type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, event storage.Event) error
	ListEventsForDay(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListEventsForWeek(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListEventsForMonth(ctx context.Context, date time.Time) ([]*storage.Event, error)
	ListEventsForNotify(ctx context.Context, date time.Time) ([]*storage.Event, error)
	DeleteOldEvents(ctx context.Context, date time.Time) error
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
		ID:        id,
		Title:     title,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(time.Hour * 24),
	}

	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, id string, domain Event) error {
	event := storage.Event{
		ID:    id,
		Title: domain.Title,
	}

	if domain.NotifyDays > 0 {
		event.NotifyDays = sql.NullInt32{Int32: domain.NotifyDays, Valid: true}
	}

	if domain.StartDate != "" {
		startDate, err := time.Parse("2006-01-02 15:04", domain.StartDate)
		if err != nil {
			return err
		}

		event.StartDate = startDate
	}

	if domain.EndDate != "" {
		endDate, err := time.Parse("2006-01-02 15:04", domain.EndDate)
		if err != nil {
			return err
		}

		event.EndDate = endDate
	}

	if domain.UserID != "" {
		event.UserID = sql.NullString{String: domain.UserID, Valid: true}
	}

	if domain.Description != "" {
		event.Description = sql.NullString{String: domain.Description, Valid: true}
	}

	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, storage.Event{ID: id})
}

func (a *App) ListEvents(ctx context.Context, date, period string) ([]Event, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	var events []*storage.Event

	switch period {
	case PeriodDay:
		events, err = a.storage.ListEventsForDay(ctx, parsedDate)
		if err != nil {
			return nil, err
		}
	case PeriodWeek:
		events, err = a.storage.ListEventsForWeek(ctx, parsedDate)
		if err != nil {
			return nil, err
		}
	case PeriodMonth:
		events, err = a.storage.ListEventsForMonth(ctx, parsedDate)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s: %w", period, ErrInvalidPeriod)
	}

	result := make([]Event, 0, len(events))
	for _, event := range events {
		var (
			description, userID string
			notifyDays          int32
		)

		if event.Description.Valid {
			description = event.Description.String
		}

		if event.UserID.Valid {
			userID = event.UserID.String
		}

		if event.NotifyDays.Valid {
			notifyDays = event.NotifyDays.Int32
		}

		result = append(result, Event{
			ID:          event.ID,
			Title:       event.Title,
			StartDate:   event.StartDate.Format("2006-01-02 15:04"),
			EndDate:     event.EndDate.Format("2006-01-02 15:04"),
			Description: description,
			UserID:      userID,
			NotifyDays:  notifyDays,
		})
	}

	return result, nil
}
