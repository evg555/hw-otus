package internalhttp

type EventsResponse struct {
	Response
	Events []Event `json:"events"`
}

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	NotifyDays  int32  `json:"notify_days"`
}

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
