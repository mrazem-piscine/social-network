package handlers

import (
	"log"
	"net/http"
	"social-network/internal/middlewares"
	ws "social-network/internal/websocket"
	"strconv"

	"github.com/gorilla/websocket"
)

var groupChatUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var groupChatManager = ws.NewGroupChatManager()

// WebSocketGroupChatHandler handles WebSocket connections for group chat
func WebSocketGroupChatHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.Atoi(r.URL.Query().Get("group_id"))
	if err != nil || groupID == 0 {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	conn, err := groupChatUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	// Register user for group chat
	groupChatManager.HandleGroupChatConnection(conn, userID, groupID)

	log.Printf("✅ User %d joined Group %d chat", userID, groupID)
}
