package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/websocket"

	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	ws "social-network/internal/websocket" // ✅ Use alias "ws"
)


// ✅ Create a Global Instance of GroupChatManager
var groupChatManager = ws.NewGroupChatManager()

// ✅ WebSocket Upgrader (Allows Cross-Origin Requests)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}


// WebSocketGroupChatHandler handles WebSocket connections for group chat
func WebSocketGroupChatHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract user ID from session or query param
	userID := middlewares.GetUserIDFromSession(r)

	// If session authentication fails, use query params
	if userID == 0 {
		userIDParam := r.URL.Query().Get("user_id")
		parsedID, err := strconv.Atoi(userIDParam)
		if err != nil {
			log.Println("❌ Failed to extract userID from query parameters")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userID = parsedID

		// ✅ Only log this if the session failed
		log.Printf("⚠️ No session found, but extracted userID %d from query parameters", userID)
	} else {
		log.Printf("✅ Verified userID %d from session", userID)
	}

	// Get Group ID
	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil || groupID == 0 {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade failed:", err)
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	// ✅ Use GroupChatManager to add user
	groupChatManager.JoinGroupChat(conn, groupID, userID)
}


func GetGroupChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
    groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
    if err != nil || groupID == 0 {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    db := config.GetDB()
	rows, err := db.Query(`SELECT id, group_id, sender_id, content, sent_at FROM group_chat_messages WHERE group_id = ? ORDER BY sent_at ASC`, groupID)
    if err != nil {
        log.Println("❌ Error fetching chat history:", err)
        http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var messages []models.GroupChatMessage
	for rows.Next() {
		var message models.GroupChatMessage
		err := rows.Scan(&message.ID, &message.GroupID, &message.SenderID, &message.Content, &message.SentAt)
		if err != nil {
			log.Println("❌ Error scanning row:", err)
			continue
		}
		messages = append(messages, message)
	}
	

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(messages)
}
