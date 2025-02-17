package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// GroupChatManager handles WebSocket connections for group chat
type GroupChatManager struct {
	GroupClients map[int]map[int]*WebSocketConn // groupID -> userID -> connection
	Mutex        sync.Mutex
}

// ✅ NewGroupChatManager initializes a new group chat manager
func NewGroupChatManager() *GroupChatManager {
	return &GroupChatManager{
		GroupClients: make(map[int]map[int]*WebSocketConn),
	}
}

// ✅ JoinGroupChat (Handles Connections)
func (gm *GroupChatManager) JoinGroupChat(conn *websocket.Conn, groupID, userID int) {
	log.Printf("📌 WebSocket connection attempt - Group ID: %d, User ID: %d", groupID, userID)

	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	// Ensure group chat exists
	if gm.GroupClients[groupID] == nil {
		gm.GroupClients[groupID] = make(map[int]*WebSocketConn)
	}

	// Close old connection if user reconnects
	if oldConn, exists := gm.GroupClients[groupID][userID]; exists {
		oldConn.Conn.Close()
		log.Printf("⚠️ User %d reconnected. Closing previous connection.", userID)
	}

	// Register user in group chat
	gm.GroupClients[groupID][userID] = &WebSocketConn{Conn: conn}
	log.Printf("✅ User %d successfully joined Group %d chat.", userID, groupID)

	// Start listening for messages in a goroutine
	go gm.ListenForMessages(conn, groupID, userID)
}

// ✅ Listen for Incoming Messages from a User
func (gm *GroupChatManager) ListenForMessages(conn *websocket.Conn, groupID, userID int) {
	for {
		var message struct {
			Content string `json:"content"`
		}
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("❌ User %d left Group %d chat. Closing connection.", userID, groupID)
			gm.RemoveUserFromGroup(groupID, userID)
			break
		}

		// Validate message content
		if message.Content == "" {
			log.Printf("⚠️ Empty message received from User %d in Group %d. Ignoring.", userID, groupID)
			continue
		}

		log.Printf("📩 Group %d | User %d: %s", groupID, userID, message.Content)

		// ✅ Broadcast the message to all users in the group
		gm.BroadcastGroupMessage(groupID, userID, message.Content)
	}
}

// ✅ Send a Message to Everyone in the Group Except Sender
func (gm *GroupChatManager) BroadcastGroupMessage(groupID, senderID int, content string) {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	for userID, conn := range gm.GroupClients[groupID] {
		if userID != senderID { // Don't send message to the sender
			conn.Mutex.Lock()
			err := conn.Conn.WriteJSON(map[string]interface{}{
				"type":      "group_message",
				"sender_id": senderID,
				"content":   content,
			})
			conn.Mutex.Unlock()

			if err != nil {
				log.Printf("❌ Error sending group message to User %d: %v", userID, err)
				conn.Conn.Close()
				delete(gm.GroupClients[groupID], userID)
			}
		}
	}
}

// ✅ Remove User From Group When They Disconnect
func (gm *GroupChatManager) RemoveUserFromGroup(groupID, userID int) {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	if conn, exists := gm.GroupClients[groupID][userID]; exists {
		conn.Conn.Close()
		delete(gm.GroupClients[groupID], userID)
		log.Printf("⚠️ User %d removed from Group %d chat.", userID, groupID)
	}
}
