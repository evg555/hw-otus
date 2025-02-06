package app

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/app/mocks"
	"github.com/evg555/hw-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	mockLogger := new(mocks.Logger)
	mockStorage := new(mocks.Storage)

	ctx := context.Background()

	mockStorage.On("CreateEvent", ctx, mock.Anything).Return(nil)

	app := New(mockLogger, mockStorage)

	err := app.CreateEvent(ctx, "test uuid", "test title")
	require.Nil(t, err)
}

func TestUpdateEvent(t *testing.T) {
	ctx := context.Background()

	type args struct {
		event Event
	}

	tests := []struct {
		name          string
		args          args
		mockFunc      func(mock *mocks.Storage)
		wantErr       bool
		expectedError error
	}{
		{
			name: "update full data in event",
			mockFunc: func(mock *mocks.Storage) {
				mock.On("UpdateEvent", ctx, "test uuid", storage.Event{
					ID:        "test uuid",
					Title:     "test title",
					StartDate: time.Date(2025, 2, 1, 9, 0, 0, 0, time.UTC),
					EndDate:   time.Date(2025, 2, 1, 10, 0, 0, 0, time.UTC),
					Description: sql.NullString{
						String: "test description",
						Valid:  true,
					},
					UserID: sql.NullString{
						String: "test user id",
						Valid:  true,
					},
					NotifyDays: sql.NullInt32{Int32: 1, Valid: true},
				}).Return(nil)
			},
			args: args{
				event: Event{
					Title:       "test title",
					StartDate:   "2025-02-01 9:00",
					EndDate:     "2025-02-01 10:00",
					Description: "test description",
					UserID:      "test user id",
					NotifyDays:  1,
				},
			},
		},
		{
			name: "event is empty",
			mockFunc: func(mock *mocks.Storage) {
				mock.On("UpdateEvent", ctx, "test uuid", storage.Event{
					ID: "test uuid",
				}).Return(nil)
			},
			args: args{
				event: Event{},
			},
		},
		{
			name:     "wrong start date",
			mockFunc: func(_ *mocks.Storage) {},
			args: args{
				event: Event{
					StartDate: "wrong date",
				},
			},
			wantErr:       true,
			expectedError: &time.ParseError{},
		},
		{
			name:     "wrong end date",
			mockFunc: func(_ *mocks.Storage) {},
			args: args{
				event: Event{
					EndDate: "wrong date",
				},
			},
			wantErr:       true,
			expectedError: &time.ParseError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := new(mocks.Logger)
			mockStorage := new(mocks.Storage)

			tt.mockFunc(mockStorage)

			app := New(mockLogger, mockStorage)
			err := app.UpdateEvent(ctx, "test uuid", tt.args.event)

			if tt.wantErr {
				require.Error(t, err)
				require.ErrorAs(t, err, &tt.expectedError)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	mockLogger := new(mocks.Logger)
	mockStorage := new(mocks.Storage)

	ctx := context.Background()

	mockStorage.On("DeleteEvent", ctx, storage.Event{ID: "test uuid"}).Return(nil)

	app := New(mockLogger, mockStorage)

	err := app.DeleteEvent(ctx, "test uuid")
	require.Nil(t, err)
}

func TestListEvents(t *testing.T) {
	ctx := context.Background()

	testErr := errors.New("test error")

	type args struct {
		date   string
		period string
	}

	tests := []struct {
		name          string
		args          args
		mockFunc      func(mock *mocks.Storage)
		wantErr       bool
		expectedError error
	}{
		{
			name: "List events for day",
			mockFunc: func(mock *mocks.Storage) {
				mock.On("ListEventsForDay", ctx,
					time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
				).Return([]*storage.Event{
					{
						ID: "test uuid", Title: "test title",
						Description: sql.NullString{String: "test description", Valid: true},
						UserID:      sql.NullString{String: "test user id", Valid: true},
					},
				}, nil)
			},
			args: args{
				date:   "2025-02-01",
				period: "day",
			},
		},
		{
			name: "List events for week",
			mockFunc: func(mock *mocks.Storage) {
				mock.On("ListEventsForWeek", ctx,
					time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
				).Return([]*storage.Event{
					{ID: "test uuid", Title: "test title"},
				}, nil)
			},
			args: args{
				date:   "2025-02-01",
				period: "week",
			},
		},
		{
			name: "List events for month",
			mockFunc: func(mock *mocks.Storage) {
				mock.On("ListEventsForMonth", ctx,
					time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
				).Return([]*storage.Event{
					{ID: "test uuid", Title: "test title"},
				}, nil)
			},
			args: args{
				date:   "2025-02-01",
				period: "month",
			},
		},
		{
			name: "storage error",
			mockFunc: func(mock *mocks.Storage) {
				mock.On("ListEventsForDay", ctx,
					time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
				).Return([]*storage.Event{}, testErr)
			},
			args: args{
				date:   "2025-02-01",
				period: "day",
			},
			wantErr:       true,
			expectedError: testErr,
		},
		{
			name:     "invalid date",
			mockFunc: func(_ *mocks.Storage) {},
			args: args{
				date:   "wrong date",
				period: "day",
			},
			wantErr:       true,
			expectedError: &time.ParseError{},
		},
		{
			name:     "invalid period",
			mockFunc: func(_ *mocks.Storage) {},
			args: args{
				date:   "2025-02-01",
				period: "wrong period",
			},
			wantErr:       true,
			expectedError: ErrInvalidPeriod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := new(mocks.Logger)
			mockStorage := new(mocks.Storage)

			tt.mockFunc(mockStorage)

			app := New(mockLogger, mockStorage)
			events, err := app.ListEvents(ctx, tt.args.date, tt.args.period)

			if tt.wantErr {
				require.Error(t, err)
				require.ErrorAs(t, err, &tt.expectedError)
				require.Empty(t, events)
			} else {
				require.Nil(t, err)
				require.Len(t, events, 1)
			}
		})
	}
}
