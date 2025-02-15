package websocket

import (
	"log"
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

// ✅ Send Notification to a User (Only if Online)
func SendNotification(userID int, message string) {
	notification := Notification{
		Type:    "event_rsvp",
		UserID:  userID,
		Message: message,
	}

	NotificationManager.Mutex.Lock()
	client, exists := NotificationManager.Clients[userID]
	NotificationManager.Mutex.Unlock()

	if exists && client != nil {
		client.Mutex.Lock()
		defer client.Mutex.Unlock()

		err := client.Conn.WriteJSON(notification)
		if err != nil {
			log.Printf("❌ Error sending WebSocket notification to User %d: %v", userID, err)
			NotificationManager.RemoveClient(userID)
		} else {
			log.Printf("✅ Sent WebSocket notification to User %d: %s", userID, message)
		}
	} else {
		log.Printf("📌 User %d is offline. Notification stored.", userID)
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
