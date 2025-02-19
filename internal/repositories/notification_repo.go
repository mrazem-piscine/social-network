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

// ✅ Save a new notification in the database
func (repo *NotificationRepository) CreateNotification(userID int, notifType, message string) error {
	log.Printf("📌 Saving notification for User %d: %s", userID, message)

	_, err := repo.DB.Exec(`
        INSERT INTO notifications (user_id, type, message, is_read, created_at) 
        VALUES (?, ?, ?, 0, CURRENT_TIMESTAMP)`, userID, notifType, message)
	if err != nil {
		log.Printf("❌ Failed to save notification for User %d: %v", userID, err)
		return err
	}

	log.Printf("✅ Notification saved for User %d", userID)
	return nil
}

// ✅ Get all notifications for a user
func (repo *NotificationRepository) GetNotifications(userID int) ([]models.Notification, error) {
	rows, err := repo.DB.Query(`
		SELECT id, user_id, type, message, is_read, created_at 
		FROM notifications 
		WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Printf("❌ Error retrieving notifications for User %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Message, &notif.IsRead, &notif.CreatedAt); err != nil {
			log.Println("❌ Error scanning notification row:", err)
			continue
		}
		notifications = append(notifications, notif)
	}

	log.Printf("📩 Retrieved %d notifications for User %d", len(notifications), userID)
	return notifications, nil
}

// ✅ Get unread notifications for a user
func (repo *NotificationRepository) GetUnreadNotifications(userID int) ([]models.Notification, error) {
	rows, err := repo.DB.Query(`
		SELECT id, user_id, type, message, is_read, created_at 
		FROM notifications 
		WHERE user_id = ? AND is_read = 0 ORDER BY created_at DESC`, userID)
	if err != nil {
		log.Printf("❌ Error retrieving unread notifications for User %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Message, &notif.IsRead, &notif.CreatedAt); err != nil {
			log.Println("❌ Error scanning unread notification row:", err)
			continue
		}
		notifications = append(notifications, notif)
	}

	log.Printf("📩 Retrieved %d unread notifications for User %d", len(notifications), userID)
	return notifications, nil
}

// ✅ Mark all notifications as read for a user
func (repo *NotificationRepository) MarkNotificationsAsRead(userID int) error {
	_, err := repo.DB.Exec(`
		UPDATE notifications SET is_read = 1 WHERE user_id = ?`, userID)
	if err != nil {
		log.Printf("❌ Error marking notifications as read for User %d: %v", userID, err)
		return err
	}

	log.Printf("✅ Marked all notifications as read for User %d", userID)
	return nil
}
