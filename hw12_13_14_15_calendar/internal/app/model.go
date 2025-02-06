package app

type Event struct {
	ID          string
	Title       string
	StartDate   string
	EndDate     string
	Description string
	UserID      string
	NotifyDays  int32
}

type Notification struct {
	EventID string
	Title   string
	Date    string
	UserID  string
}
