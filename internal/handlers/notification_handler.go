package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
	ws "social-network/internal/websocket" // alias for our internal websocket package

	"github.com/gorilla/websocket"
)

// WebSocketNotificationHandler handles real-time notifications via WebSocket.
func WebSocketNotificationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user_id from query parameters.
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket.
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("❌ Failed to upgrade WebSocket:", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	// Set read deadline and pong handler to keep connection alive.
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Register the client connection for notifications.
	ws.NotificationManager.RegisterClient(userID, conn)
	log.Printf("✅ WebSocket connected for User %d", userID)

	// Start a ping ticker to send ping messages every 30 seconds.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Run ping in a separate goroutine.
	go func() {
		for {
			<-ticker.C
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("❌ Failed to send ping to User %d: %v", userID, err)
				return
			}
		}
	}()

	// Keep the connection alive by reading messages.
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("❌ WebSocket read error for User %d: %v", userID, err)
			ws.NotificationManager.RemoveClient(userID)
			break
		}
		// You can handle incoming messages here if needed.
	}
} // GetNotificationsHandler fetches all notifications for the authenticated user.
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
		log.Println("❌ Error fetching notifications:", err)
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	// Optionally, mark them as read immediately (or do this in a separate endpoint)
	_, err = db.Exec("UPDATE notifications SET is_read = 1 WHERE user_id = ?", userID)
	if err != nil {
		log.Println("❌ Error marking notifications as read:", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifications)
}

// SendNotificationHandler allows an API request to create (and send) a notification.
func SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON request body
	var requestBody struct {
		UserID  int    `json:"user_id"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Get database connection and initialize repository
	db := config.GetDB()
	repo := repositories.NewNotificationRepository(db)

	// Save the notification in the database
	err := repo.CreateNotification(requestBody.UserID, requestBody.Type, requestBody.Message)
	if err != nil {
		log.Println("❌ Failed to save notification:", err)
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	log.Printf("📩 Notification sent to User %d: %s", requestBody.UserID, requestBody.Message)

	// Also, if the user is online, send it via WebSocket
	ws.SendNotification(requestBody.UserID, requestBody.Type, requestBody.Message)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification sent successfully"})
}

// MarkNotificationsAsReadHandler marks all notifications as read for the authenticated user.
func MarkNotificationsAsReadHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	db := config.GetDB()
	repo := repositories.NewNotificationRepository(db)

	err := repo.MarkNotificationsAsRead(userID)
	if err != nil {
		http.Error(w, "Failed to mark notifications as read", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "All notifications marked as read"})
}
