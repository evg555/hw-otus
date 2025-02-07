package storage

import (
	"database/sql"
	"time"
)

type Event struct {
	ID          string         `db:"uuid"`
	Title       string         `db:"title"`
	StartDate   time.Time      `db:"start_date"`
	EndDate     time.Time      `db:"end_date"`
	Description sql.NullString `db:"description"`
	UserID      sql.NullString `db:"user_id"`
	NotifyDays  sql.NullInt32  `db:"notify_days"`
}
