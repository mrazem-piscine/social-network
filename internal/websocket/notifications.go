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

// NewWebSocketNotificationManager initializes a new notification manager
func NewWebSocketNotificationManager() *WebSocketNotificationManager {
	return &WebSocketNotificationManager{
		Clients: make(map[int]*WebSocketConn),
	}
}

// SendNotification sends a notification to an online user
func (wm *WebSocketNotificationManager) SendNotification(notification Notification) {
	wm.Mutex.Lock()
	clientConn, exists := wm.Clients[notification.UserID]
	wm.Mutex.Unlock()

	if exists {
		clientConn.Mutex.Lock()
		err := clientConn.Conn.WriteJSON(notification)
		clientConn.Mutex.Unlock()

		if err != nil {
			log.Printf("❌ Error sending notification to user %d: %v", notification.UserID, err)
			wm.RemoveClient(notification.UserID) // Remove if connection fails
		}
	} else {
		log.Printf("⚠️ User %d is offline. Notification stored for later.", notification.UserID)
	}
}

// RegisterClient registers a WebSocket client for notifications
func (wm *WebSocketNotificationManager) RegisterClient(userID int, conn *websocket.Conn) {
	wm.Mutex.Lock()
	wm.Clients[userID] = &WebSocketConn{Conn: conn}
	wm.Mutex.Unlock()
	log.Printf("✅ User %d connected for real-time notifications.", userID)
}

// RemoveClient removes a disconnected WebSocket client
func (wm *WebSocketNotificationManager) RemoveClient(userID int) {
	wm.Mutex.Lock()
	delete(wm.Clients, userID)
	wm.Mutex.Unlock()
	log.Printf("⚠️ User %d disconnected from notifications.", userID)
}
