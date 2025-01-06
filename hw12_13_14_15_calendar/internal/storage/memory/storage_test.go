package memorystorage

import (
	"context"
	"testing"
	"time"

	internalstorage "github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		ctx := context.Background()

		storage := New()
		event := internalstorage.Event{ID: "1"}

		err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)

		err = storage.CreateEvent(ctx, event)
		require.Error(t, err)
		require.ErrorIs(t, err, internalstorage.ErrEventAlreadyExists)

		event.Title = "test"
		err = storage.UpdateEvent(ctx, "1", event)
		require.NoError(t, err)

		event.Title = "test"
		err = storage.UpdateEvent(ctx, "2", event)
		require.Error(t, err)
		require.ErrorIs(t, err, internalstorage.ErrEventNotExists)

		event2 := internalstorage.Event{ID: "2"}
		err = storage.CreateEvent(ctx, event2)
		require.NoError(t, err)

		err = storage.DeleteEvent(ctx, event)
		require.NoError(t, err)

		err = storage.DeleteEvent(ctx, event)
		require.NoError(t, err)
	})

	t.Run("list events", func(t *testing.T) {
		ctx := context.Background()

		date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		storage := New()
		events := []internalstorage.Event{
			{ID: "1", StartDate: time.Date(2020, 1, 1, 13, 0, 0, 0, time.UTC)},
			{ID: "2", StartDate: time.Date(2020, 1, 1, 8, 0, 0, 0, time.UTC)},
			{ID: "3", StartDate: time.Date(2020, 1, 2, 9, 0, 0, 0, time.UTC)},
			{ID: "4", StartDate: time.Date(2020, 1, 15, 22, 0, 0, 0, time.UTC)},
			{ID: "5", StartDate: time.Date(2020, 2, 1, 11, 0, 0, 0, time.UTC)},
		}

		for _, event := range events {
			err := storage.CreateEvent(ctx, event)
			require.NoError(t, err)
		}

		got, err := storage.ListEventsForDay(ctx, date)
		require.NoError(t, err)
		require.Len(t, got, 2)

		got, err = storage.ListEventsForWeek(ctx, date)
		require.NoError(t, err)
		require.Len(t, got, 3)

		got, err = storage.ListEventsForMonth(ctx, date)
		require.NoError(t, err)
		require.Len(t, got, 4)
	})
}
