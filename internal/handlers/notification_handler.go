package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
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
	repo := repositories.NewNotificationRepository(db)

	notifications, err := repo.GetNotifications(userID)
	if err != nil {
		log.Println("Error getting notifications:", err)
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifications)
}
