package websocket

import (
	"log"
	"sync"

	"social-network/internal/config"
	"social-network/internal/models"

	"github.com/gorilla/websocket"
)

// WebSocketConn wraps a WebSocket connection.
type WebSocketConn struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

// WebSocketNotificationManager manages WebSocket notifications.
type WebSocketNotificationManager struct {
	Clients map[int]*WebSocketConn
	Mutex   sync.Mutex
}

// Notification represents a real-time notification (for WebSocket output).
type Notification struct {
	Type    string `json:"type"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

// Global Notification Manager instance.
var NotificationManager = &WebSocketNotificationManager{
	Clients: make(map[int]*WebSocketConn),
}

// NewWebSocketNotificationManager initializes a new notification manager.
func NewWebSocketNotificationManager() *WebSocketNotificationManager {
	return &WebSocketNotificationManager{
		Clients: make(map[int]*WebSocketConn),
	}
}

// SendNotification sends a notification to a user via WebSocket.
// If sending fails or the user is offline, it stores the notification in the database.
func SendNotification(userID int, notifType, message string) {
	notification := models.Notification{
		UserID:  userID,
		Type:    notifType,
		Message: message,
		IsRead:  false,
	}

	NotificationManager.Mutex.Lock()
	client, exists := NotificationManager.Clients[userID]
	NotificationManager.Mutex.Unlock()

	if exists && client != nil {
		client.Mutex.Lock()
		err := client.Conn.WriteJSON(notification)
		client.Mutex.Unlock()

		if err != nil {
			log.Printf("❌ Error sending WebSocket notification to User %d: %v", userID, err)
			NotificationManager.RemoveClient(userID)
			storeNotification(notification) // Fallback: store in DB
		} else {
			log.Printf("✅ WebSocket notification sent to User %d", userID)
		}
	} else {
		log.Printf("📌 User %d is offline. Storing notification.", userID)
		// storeNotification(notification)
	}
}

// storeNotification stores the notification in the database (fallback when WebSocket sending fails).
func storeNotification(notification models.Notification) {
	db := config.GetDB()
	_, err := db.Exec(`
		INSERT INTO notifications (user_id, type, message, is_read, created_at) 
		VALUES (?, ?, ?, 0, CURRENT_TIMESTAMP)`,
		notification.UserID, notification.Type, notification.Message)
	if err != nil {
		log.Printf("❌ Failed to store notification for User %d: %v", notification.UserID, err)
	} else {
		log.Printf("✅ Notification stored for User %d", notification.UserID)
	}
}

// RegisterClient registers a WebSocket client for notifications.
func (wm *WebSocketNotificationManager) RegisterClient(userID int, conn *websocket.Conn) {
	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	// If a connection already exists, close it.
	if oldClient, exists := wm.Clients[userID]; exists {
		log.Printf("⚠️ Closing previous notification connection for User %d.", userID)
		oldClient.Conn.Close()
		delete(wm.Clients, userID)
	}

	wm.Clients[userID] = &WebSocketConn{Conn: conn}
	log.Printf("✅ User %d connected for real-time notifications.", userID)
}

// RemoveClient removes a disconnected WebSocket client.
func (wm *WebSocketNotificationManager) RemoveClient(userID int) {
	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	if client, exists := wm.Clients[userID]; exists {
		client.Conn.Close()
		delete(wm.Clients, userID)
		log.Printf("⚠️ User %d disconnected from notifications.", userID)
	}
}
