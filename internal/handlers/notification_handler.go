package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
	ws "social-network/internal/websocket" // ✅ Assign alias to avoid conflicts

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var notificationManager = ws.NewWebSocketNotificationManager() // ✅ Use alias `ws`

// WebSocketNotificationHandler handles real-time notifications
func WebSocketNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	// Register user for real-time notifications
	notificationManager.RegisterClient(userID, conn)

	log.Printf("✅ User %d connected for real-time notifications", userID)
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
