package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	m  map[string]storage.Event
	mu sync.RWMutex
}

func New() *Storage {
	return &Storage{
		m: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[event.ID]; ok {
		return storage.ErrEventAlreadyExists
	}

	s.m[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[id]; !ok {
		return storage.ErrEventNotExists
	}

	s.m[id] = event

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	delete(s.m, event.ID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) ListEventsForDay(_ context.Context, date time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	dateStr := date.Format("2006-01-02")

	s.mu.RLock()
	for _, event := range s.m {
		eventDate := event.StartDate.Format("2006-01-02")

		if eventDate == dateStr {
			events = append(events, &event)
		}
	}
	s.mu.RUnlock()

	return events, nil
}

func (s *Storage) ListEventsForWeek(_ context.Context, date time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	startOfWeek := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	s.mu.RLock()
	for _, event := range s.m {
		if event.StartDate.After(startOfWeek.Add(-1*time.Second)) && event.StartDate.Before(endOfWeek.Add(24*time.Hour)) {
			events = append(events, &event)
		}
	}
	s.mu.RUnlock()

	return events, nil
}

func (s *Storage) ListEventsForMonth(_ context.Context, date time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	s.mu.RLock()
	for _, event := range s.m {
		if event.StartDate.After(startOfMonth.Add(-1*time.Second)) && event.StartDate.Before(endOfMonth.Add(24*time.Hour)) {
			events = append(events, &event)
		}
	}
	s.mu.RUnlock()

	return events, nil
}

func (s *Storage) ListEventsForNotify(_ context.Context, date time.Time) ([]*storage.Event, error) {
	var events []*storage.Event

	s.mu.RLock()
	for _, event := range s.m {
		if event.NotifyDays.Int32 == 0 {
			continue
		}

		needNotifyDate := date.Add(time.Duration(event.NotifyDays.Int32) * 24 * time.Hour).Format("2006-01-02")

		if needNotifyDate == event.StartDate.Format("2006-01-02") {
			events = append(events, &event)
		}
	}
	s.mu.RUnlock()

	return events, nil
}

func (s *Storage) DeleteOldEvents(_ context.Context, date time.Time) error {
	oldDate := date.Add(-(365 * 24 * time.Hour))

	s.mu.Lock()
	for _, event := range s.m {
		if event.EndDate.Before(oldDate) {
			delete(s.m, event.ID)
		}
	}
	s.mu.Unlock()

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}
