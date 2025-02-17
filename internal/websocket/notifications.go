package websocket

import (
	"log"
	"social-network/internal/config"
	"social-network/internal/models"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketConn wraps a WebSocket connection
type WebSocketConn struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

// WebSocketNotificationManager manages WebSocket notifications
type WebSocketNotificationManager struct {
	Clients map[int]*WebSocketConn
	Mutex   sync.Mutex
}

// Notification represents a real-time notification
type Notification struct {
	Type    string `json:"type"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

// ✅ Global Notification Manager (Pointer)
var NotificationManager = &WebSocketNotificationManager{
	Clients: make(map[int]*WebSocketConn),
}

// NewWebSocketNotificationManager initializes a new notification manager
func NewWebSocketNotificationManager() *WebSocketNotificationManager {
	return &WebSocketNotificationManager{
		Clients: make(map[int]*WebSocketConn),
	}
}

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
			log.Printf("❌ Error sending WebSocket notification: %v", err)
			NotificationManager.RemoveClient(userID)
			storeNotification(notification) // Store in DB if sending fails
		}
	} else {
		storeNotification(notification) // Store in DB if offline
		log.Printf("📌 User %d is offline. Notification stored.", userID)
	}
}

// **Function to store notifications in the database**
func storeNotification(notification models.Notification) {
	db := config.GetDB()
	_, err := db.Exec(`
		INSERT INTO notifications (user_id, type, message, is_read, created_at) 
		VALUES (?, ?, ?, 0, CURRENT_TIMESTAMP)`,
		notification.UserID, notification.Type, notification.Message)

	if err != nil {
		log.Printf("❌ Failed to store notification: %v", err)
	}
}

// ✅ Register a WebSocket Client for Notifications
func (wm *WebSocketNotificationManager) RegisterClient(userID int, conn *websocket.Conn) {
	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	// ✅ Close previous connection if user reconnects
	if oldClient, exists := wm.Clients[userID]; exists {
		log.Printf("⚠️ Closing old connection for User %d.", userID)
		oldClient.Conn.Close()
		delete(wm.Clients, userID)
	}

	wm.Clients[userID] = &WebSocketConn{Conn: conn}
	log.Printf("✅ User %d connected for real-time notifications.", userID)
}

// ✅ Remove a Disconnected WebSocket Client
func (wm *WebSocketNotificationManager) RemoveClient(userID int) {
	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	if client, exists := wm.Clients[userID]; exists {
		client.Conn.Close()
		delete(wm.Clients, userID)
		log.Printf("⚠️ User %d disconnected from notifications.", userID)
	}
}
