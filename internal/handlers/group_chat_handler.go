package handlers

import (
	"log"
	"net/http"
	"social-network/internal/middlewares"
	ws "social-network/internal/websocket"
	"strconv"

	"github.com/gorilla/websocket"
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
