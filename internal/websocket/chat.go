package websocket

import (
	"database/sql"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketManager manages WebSocket connections
type WebSocketManager struct {
	Clients map[string]*WebSocketConn
	Mutex   sync.Mutex
}

// WebSocketConn represents a WebSocket connection


// NewWebSocketManager initializes a new WebSocketManager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients: make(map[string]*WebSocketConn),
	}
}

// HandleConnection manages WebSocket connections
func (wm *WebSocketManager) HandleConnection(conn *websocket.Conn, username string, db *sql.DB) {
	wsConn := &WebSocketConn{Conn: conn}
	wm.Mutex.Lock()
	wm.Clients[username] = wsConn
	wm.Mutex.Unlock()

	defer func() {
		wm.Mutex.Lock()
		delete(wm.Clients, username)
		wm.Mutex.Unlock()
		wsConn.Conn.Close()
	}()

	for {
		var msg struct {
			Type     string `json:"type"`
			Sender   string `json:"sender"`
			Receiver string `json:"receiver"`
			Content  string `json:"content"`
			Time     string `json:"time"`
		}

		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from %s: %v", username, err)
			break
		}

		// ✅ Keep the connection alive by sending an acknowledgment
		response := struct {
			Status string `json:"status"`
		}{
			Status: "Received",
		}
		conn.WriteJSON(response)
	}
}
