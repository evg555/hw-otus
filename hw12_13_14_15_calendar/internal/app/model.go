package app

import "time"

type Event struct {
	ID               string
	Title            string
	StartDate        string
	EndDate          string
	Description      string
	UserID           string
	NotificationTime time.Duration
}
