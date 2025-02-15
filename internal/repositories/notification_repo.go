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

// CreateNotification inserts a new notification into the database
func (repo *NotificationRepository) CreateNotification(userID int, message string) error {
	_, err := repo.DB.Exec(`
        INSERT INTO notifications (user_id, message, is_read, created_at) 
        VALUES (?, ?, 0, CURRENT_TIMESTAMP)`, userID, message)

	if err != nil {
		log.Println("❌ Error inserting notification:", err)
	}
	return err
}

// GetNotifications fetches notifications for a user
func (repo *NotificationRepository) GetNotifications(userID int) ([]models.Notification, error) {
	rows, err := repo.DB.Query("SELECT id, user_id, type, message, created_at FROM notifications WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error retrieving notifications:", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Message, &notif.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// GetUnreadNotifications fetches unread notifications for a user
func (repo *NotificationRepository) GetUnreadNotifications(userID int) ([]models.Notification, error) {
	rows, err := repo.DB.Query("SELECT id, user_id, type, message, created_at FROM notifications WHERE user_id = ? AND is_read = 0", userID)
	if err != nil {
		log.Println("Error retrieving unread notifications:", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Message, &notif.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// MarkNotificationsAsRead marks all notifications as read for a user
func (repo *NotificationRepository) MarkNotificationsAsRead(userID int) error {
	_, err := repo.DB.Exec("UPDATE notifications SET is_read = 1 WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error marking notifications as read:", err)
		return err
	}
	return nil
}
