package websocket

import (
	"log"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

// GroupChatManager handles WebSocket connections for group chat
type GroupChatManager struct {
	GroupClients map[int]map[int]*WebSocketConn // groupID -> userID -> connection
	Mutex        sync.Mutex
}

// NewGroupChatManager initializes a new group chat manager
func NewGroupChatManager() *GroupChatManager {
	return &GroupChatManager{
		GroupClients: make(map[int]map[int]*WebSocketConn),
	}
}

// HandleGroupChatConnection manages WebSocket connections for group chat
func (gm *GroupChatManager) HandleGroupChatConnection(conn *websocket.Conn, userID, groupID int) {
	gm.Mutex.Lock()
	if gm.GroupClients[groupID] == nil {
		gm.GroupClients[groupID] = make(map[int]*WebSocketConn)
	}
	gm.GroupClients[groupID][userID] = &WebSocketConn{Conn: conn}
	gm.Mutex.Unlock()

	log.Printf("✅ User %d joined Group %d chat.", userID, groupID)

	// Listen for messages
	for {
		var message map[string]string
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("❌ User %d left Group %d chat.", userID, groupID)
			gm.Mutex.Lock()
			delete(gm.GroupClients[groupID], userID)
			gm.Mutex.Unlock()
			break
		}

		content := message["content"]
		log.Printf("📩 Group %d | User %d: %s", groupID, userID, content)

		// Send the message to all users in the group
		gm.BroadcastGroupMessage(groupID, userID, content)
	}
}

// BroadcastGroupMessage sends a message to all users in a group chat
func (gm *GroupChatManager) BroadcastGroupMessage(groupID, senderID int, content string) {
	gm.Mutex.Lock()
	defer gm.Mutex.Unlock()

	for userID, conn := range gm.GroupClients[groupID] {
		if userID != senderID {
			conn.Mutex.Lock()
			err := conn.Conn.WriteJSON(map[string]string{
				"sender_id": strconv.Itoa(senderID),
				"content":   content,
			})
			conn.Mutex.Unlock()

			if err != nil {
				log.Printf("❌ Error sending group message to User %d: %v", userID, err)
			}
		}
	}
}
