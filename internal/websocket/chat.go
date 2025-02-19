package websocket

import (
	"log"
	"strconv"
	"sync"

	"social-network/internal/config"
	"social-network/internal/repositories"

	"github.com/gorilla/websocket"
)

// ChatManager handles WebSocket chat connections
type ChatManager struct {
	Clients map[int]*WebSocketConn // userID -> WebSocket connection
	Mutex   sync.Mutex
}

// ✅ Create a new ChatManager instance
func NewChatManager() *ChatManager {
	return &ChatManager{
		Clients: make(map[int]*WebSocketConn),
	}
}

// ✅ Register user connection for WebSocket chat
func (cm *ChatManager) HandleChatConnection(conn *websocket.Conn, userID int) {
	log.Printf("📌 WebSocket connection attempt - User ID: %d", userID)

	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()

	// Close old connection if user reconnects
	if oldConn, exists := cm.Clients[userID]; exists {
		oldConn.Conn.Close()
		log.Printf("⚠️ User %d reconnected. Closing previous connection.", userID)
	}

	// Register user in the chat system
	cm.Clients[userID] = &WebSocketConn{Conn: conn}
	log.Printf("✅ User %d successfully joined private chat.", userID)

	// Start listening for messages
	go cm.ListenForMessages(conn, userID)
}

// ✅ Listen for Incoming Messages from a User
func (cm *ChatManager) ListenForMessages(conn *websocket.Conn, userID int) {
	for {
		var message struct {
			ReceiverID int    `json:"receiver_id"`
			Content    string `json:"content"`
		}
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("❌ User %d disconnected from chat. Closing connection.", userID)
			cm.RemoveClient(userID)
			break
		}

		// Validate message content
		if message.Content == "" {
			log.Printf("⚠️ Empty message received from User %d. Ignoring.", userID)
			continue
		}

		log.Printf("📩 Private Chat | User %d -> User %d: %s", userID, message.ReceiverID, message.Content)

		// ✅ Send the message to the receiver and save it
		cm.SendMessage(userID, message.ReceiverID, message.Content)
	}
}

// ✅ Send a Message to a Specific User and Save to Database
func (cm *ChatManager) SendMessage(senderID int, receiverID int, content string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()

	// ✅ Save message in the database
	db := config.GetDB()
	repo := repositories.NewChatRepository(db)
	err := repo.SaveMessage(senderID, receiverID, content)
	if err != nil {
		log.Printf("❌ Failed to save message to database: %v", err)
	}

	// ✅ Check if receiver is connected
	receiverConn, exists := cm.Clients[receiverID]
	if !exists {
		log.Printf("⚠️ User %d is NOT connected. Message will be stored in the database.", receiverID)
		return
	}

	// ✅ Send message to receiver via WebSocket
	receiverConn.Mutex.Lock()
	err = receiverConn.Conn.WriteJSON(map[string]string{
		"sender_id": strconv.Itoa(senderID),
		"content":   content,
	})
	receiverConn.Mutex.Unlock()

	if err != nil {
		log.Printf("❌ Error sending WebSocket message to User %d: %v", receiverID, err)
	}
}

// ✅ Remove User from the chat when they disconnect
func (cm *ChatManager) RemoveClient(userID int) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()

	if conn, exists := cm.Clients[userID]; exists {
		conn.Conn.Close()
		delete(cm.Clients, userID)
		log.Printf("⚠️ User %d disconnected from chat and removed.", userID)
	}
}

func (cm *ChatManager) HandleConnection(conn *websocket.Conn, userID int) {
	cm.Mutex.Lock()
	cm.Clients[userID] = &WebSocketConn{Conn: conn}
	cm.Mutex.Unlock()

	log.Printf("✅ User %d connected to WebSocket chat", userID)

	go cm.ListenForMessages(conn, userID)
}
