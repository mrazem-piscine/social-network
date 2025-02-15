package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
)

// GroupEventRepository handles group event-related operations
type GroupEventRepository struct {
	DB *sql.DB
}

// ✅ NewGroupEventRepository initializes the repository
func NewGroupEventRepository(db *sql.DB) *GroupEventRepository {
	return &GroupEventRepository{DB: db}
}

// ✅ CreateGroupEvent adds a new event to a group
func (repo *GroupEventRepository) CreateGroupEvent(event *models.GroupEvent) error {
	_, err := repo.DB.Exec(`
        INSERT INTO group_events (group_id, creator_id, title, description, event_date)
        VALUES (?, ?, ?, ?, ?)`,

		event.GroupID, event.CreatorID, event.Title, event.Description, event.EventDate)
	return err
}

// ✅ GetEventByID retrieves an event's details by ID
func (repo *GroupEventRepository) GetEventByID(eventID int) (*models.GroupEvent, error) {
	var event models.GroupEvent
	err := repo.DB.QueryRow(`
		SELECT id, group_id, creator_id, title, description, event_date 
		FROM group_events WHERE id = ?`, eventID).
		Scan(&event.ID, &event.GroupID, &event.CreatorID, &event.Title, &event.Description, &event.EventDate)

	if err != nil {
		log.Println("❌ Error fetching event:", err)
		return nil, err
	}
	return &event, nil
}

// ✅ RSVPToEvent allows users to RSVP for an event
func (repo *GroupEventRepository) RSVPToEvent(eventID, userID int, status string) error {
	_, err := repo.DB.Exec(`
	INSERT INTO event_rsvps (event_id, user_id, status, created_at) 
	VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(event_id, user_id) 
	DO UPDATE SET status = excluded.status, created_at = CURRENT_TIMESTAMP`,
		eventID, userID, status)

	if err != nil {
		log.Println("❌ Error updating RSVP:", err)
	}
	return err
}
