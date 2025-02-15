package repositories

import (
	"database/sql"
)

// EventRSVPRepository handles event RSVP operations
type EventRSVPRepository struct {
	DB *sql.DB
}

// NewEventRSVPRepository creates a new instance of EventRSVPRepository
func NewEventRSVPRepository(db *sql.DB) *EventRSVPRepository {
	return &EventRSVPRepository{DB: db}
}

// RSVPToEvent lets users mark themselves as "going" or "not going"
func (repo *EventRSVPRepository) RSVPToEvent(eventID, userID int, status string) error {
	_, err := repo.DB.Exec(`
        INSERT INTO event_rsvps (event_id, user_id, status)
        VALUES (?, ?, ?)
        ON CONFLICT(event_id, user_id) DO UPDATE SET status = excluded.status`,
		eventID, userID, status)
	return err
}

// GetRSVPCount returns the number of users who are "going" to an event
func (repo *EventRSVPRepository) GetRSVPCount(eventID int) (int, error) {
	var count int
	err := repo.DB.QueryRow(`
        SELECT COUNT(*) FROM event_rsvps WHERE event_id = ? AND status = 'going'`,
		eventID).Scan(&count)
	return count, err
}
