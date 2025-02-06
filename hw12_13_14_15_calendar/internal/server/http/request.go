package internalhttp

type CreateRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
type UpdateEventRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	NotifyDays  int32  `json:"notify_days"`
}

type ListEventsRequest struct {
	Date   string
	Period string
}
