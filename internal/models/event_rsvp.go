package models

// EventRSVP represents a user's RSVP status for an event
type EventRSVP struct {
	ID        int    `json:"id"`
	EventID   int    `json:"event_id"`
	UserID    int    `json:"user_id"`
	Status    string `json:"status"` // "going" or "not going"
	CreatedAt string `json:"created_at"`
}
