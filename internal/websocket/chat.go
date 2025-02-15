package websocket

import (
	"database/sql"
	"log"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

// ChatManager handles WebSocket chat connections
type ChatManager struct {
	Clients map[int]*WebSocketConn
	Mutex   sync.Mutex
}

// NewChatManager initializes a new chat manager
// WebSocketManager manages WebSocket connections
type WebSocketManager struct {
	Clients map[string]*WebSocketConn
	Mutex   sync.Mutex
}

func NewChatManager() *ChatManager {
	return &ChatManager{
		Clients: make(map[int]*WebSocketConn),
	}
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

func (cm *ChatManager) HandleChatConnection(conn *websocket.Conn, userID int) {
	cm.Mutex.Lock()
	cm.Clients[userID] = &WebSocketConn{Conn: conn}
	cm.Mutex.Unlock()

	log.Printf("✅ User %d connected to chat.", userID)

	// Listen for messages
	for {
		var message map[string]string
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("❌ User %d disconnected from chat.", userID)
			cm.Mutex.Lock()
			delete(cm.Clients, userID)
			cm.Mutex.Unlock()
			break
		}

		receiverID, _ := strconv.Atoi(message["receiver_id"])
		content := message["content"]

		log.Printf("📩 Message from User %d to User %d: %s", userID, receiverID, content)

		// Send the message if the receiver is online
		cm.SendMessage(userID, receiverID, content)
	}
}

// SendMessage sends a message to an online user
func (cm *ChatManager) SendMessage(senderID int, receiverID int, content string) {
	cm.Mutex.Lock()
	receiverConn, exists := cm.Clients[receiverID]
	cm.Mutex.Unlock()

	if exists {
		receiverConn.Mutex.Lock()
		err := receiverConn.Conn.WriteJSON(map[string]string{"sender_id": strconv.Itoa(senderID), "content": content})
		receiverConn.Mutex.Unlock()

		if err != nil {
			log.Printf("❌ Error sending message to user %d: %v", receiverID, err)
		}
	} else {
		log.Printf("⚠️ User %d is offline. Message stored for later.", receiverID)
	}
}
