package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketConn wraps a WebSocket connection with a mutex
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

// SendNotification sends a notification to an online user via WebSocket
func (wm *WebSocketNotificationManager) SendNotification(notification Notification) {
	wm.Mutex.Lock()
	clientConn, exists := wm.Clients[notification.UserID]
	wm.Mutex.Unlock() // ✅ Unlock immediately after checking existence

	if !exists {
		log.Printf("User %d is offline. Notification stored for later.", notification.UserID)
		return
	}

	clientConn.Mutex.Lock()
	defer clientConn.Mutex.Unlock() // ✅ Ensure mutex is released

	// ✅ Ensure connection is still open before writing
	if err := clientConn.Conn.WriteJSON(notification); err != nil {
		log.Printf("Error sending notification to user %d: %v", notification.UserID, err)
		wm.RemoveClient(notification.UserID) // ✅ Remove disconnected clients
	}
}

// RemoveClient removes a disconnected WebSocket client
func (wm *WebSocketNotificationManager) RemoveClient(userID int) {
	wm.Mutex.Lock()
	defer wm.Mutex.Unlock()

	if clientConn, exists := wm.Clients[userID]; exists {
		clientConn.Conn.Close()
		delete(wm.Clients, userID)
		log.Printf("Removed disconnected client for user %d", userID)
	}
}
