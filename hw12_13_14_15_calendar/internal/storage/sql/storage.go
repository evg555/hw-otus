package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib" // for PostgreSQL driver
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfg config.DBConf) *Storage {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		panic(fmt.Sprintf("database init error: %v", err))
	}

	return &Storage{
		db: db,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	err := s.db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	query := `INSERT INTO events (uuid, title, start_date, end_date, description, user_id, notify_days)
    			VALUES (:uuid, :title, :start_date, :end_date, :description, :user_id, :notify_days)`

	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	event.ID = id

	query := `UPDATE events SET title=:title, start_date=:start_date, end_date=:end_date, description=:description, 
                  user_id=:user_id, notify_days=:notify_days WHERE uuid = :uuid`

	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, event storage.Event) error {
	query := `DELETE FROM events WHERE uuid = $1`

	_, err := s.db.ExecContext(ctx, query, event.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	query := `SELECT uuid, title, start_date, end_date, description, user_id, notify_days
				FROM events WHERE start_date::DATE = $1`

	stmt, err := s.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var events []*storage.Event

	err = stmt.SelectContext(ctx, &events, date)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ListEventsForWeek(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	startOfWeek := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	query := `SELECT uuid, title, start_date, end_date, description, user_id, notify_days 
				FROM events WHERE start_date::DATE between $1 and $2`

	stmt, err := s.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var events []*storage.Event

	err = stmt.SelectContext(ctx, &events, startOfWeek, endOfWeek)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ListEventsForMonth(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	query := `SELECT uuid, title, start_date, end_date, description, user_id, notify_days
				FROM events WHERE start_date::DATE between $1 and $2`

	stmt, err := s.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var events []*storage.Event

	err = stmt.SelectContext(ctx, &events, startOfMonth, endOfMonth)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ListEventsForNotify(ctx context.Context, date time.Time) ([]*storage.Event, error) {
	query := `SELECT uuid, title, start_date, user_id FROM events
				WHERE start_date::DATE = ($1::TIMESTAMP + INTERVAL '1 day' * notify_days)::DATE`

	stmt, err := s.db.Preparex(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var events []*storage.Event

	err = stmt.SelectContext(ctx, &events, date)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) DeleteOldEvents(ctx context.Context, date time.Time) error {
	oneYearAgo := date.AddDate(-1, 0, 0)

	query := `DELETE FROM events WHERE end_date < $1`

	_, err := s.db.ExecContext(ctx, query, oneYearAgo)
	if err != nil {
		return err
	}

	return nil
}
