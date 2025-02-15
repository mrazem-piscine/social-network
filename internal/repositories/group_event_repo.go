package repositories

import (
	"database/sql"
	"social-network/internal/models"
)

// GroupEventRepository handles group event-related operations
type GroupEventRepository struct {
	DB *sql.DB
}

// NewGroupEventRepository creates a new instance of GroupEventRepository
func NewGroupEventRepository(db *sql.DB) *GroupEventRepository {
	return &GroupEventRepository{DB: db}
}

// CreateGroupEvent adds a new event to a group
func (repo *GroupEventRepository) CreateGroupEvent(event *models.GroupEvent) error {
	_, err := repo.DB.Exec(`
        INSERT INTO group_events (group_id, creator_id, title, description, event_date)
        VALUES (?, ?, ?, ?, ?)`,
		event.GroupID, event.CreatorID, event.Title, event.Description, event.EventDate)
	return err
}
