package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
)

// NotificationRepository handles database operations for notifications
type NotificationRepository struct {
	DB *sql.DB
}

// NewNotificationRepository creates a new instance of NotificationRepository
func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

// CreateNotification inserts a notification into the database
func (repo *NotificationRepository) CreateNotification(userID int, notifType, message string) error {
	_, err := repo.DB.Exec(`
		INSERT INTO notifications (user_id, type, message) 
		VALUES (?, ?, ?)`, userID, notifType, message)

	if err != nil {
		log.Println("Error inserting notification:", err)
	}
	return err
}

// GetNotifications retrieves unread notifications for a user
func (repo *NotificationRepository) GetNotifications(userID int) ([]models.Notification, error) {
	rows, err := repo.DB.Query(`
		SELECT id, type, message, is_read, created_at 
		FROM notifications 
		WHERE user_id = ? 
		ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Println("Error fetching notifications:", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		err := rows.Scan(&notif.ID, &notif.Type, &notif.Message, &notif.IsRead, &notif.CreatedAt)
		if err != nil {
			log.Println("Error scanning notification:", err)
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, nil
}
