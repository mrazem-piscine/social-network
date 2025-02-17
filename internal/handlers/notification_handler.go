package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	ws "social-network/internal/websocket" // ✅ Assign alias to avoid conflicts
	"strconv"

	"github.com/gorilla/websocket"
)

// WebSocketNotificationHandler handles real-time notifications
func WebSocketNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024) // ✅ Fix: Added missing buffer sizes
	if err != nil {
		log.Println("❌ Failed to upgrade WebSocket:", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	// Store connection in WebSocket manager
	ws.NotificationManager.RegisterClient(userID, conn)
	log.Printf("✅ WebSocket connected for User %d", userID)
}

// GetNotificationsHandler fetches notifications for a user
func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db := config.GetDB()
	rows, err := db.Query(`
		SELECT id, type, message, is_read, created_at 
		FROM notifications WHERE user_id = ?`, userID)
	if err != nil {
		log.Println("❌ Error fetching notifications:", err)
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		err := rows.Scan(&notif.ID, &notif.Type, &notif.Message, &notif.IsRead, &notif.CreatedAt)
		if err != nil {
			log.Println("❌ Error scanning notifications:", err)
			continue
		}
		notifications = append(notifications, notif)
	}

	// ✅ Mark notifications as read
	_, err = db.Exec("UPDATE notifications SET is_read = 1 WHERE user_id = ?", userID)
	if err != nil {
		log.Println("❌ Error marking notifications as read:", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifications)
}
