package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	"social-network/internal/repositories"
	ws "social-network/internal/websocket"
	"strconv"

	"github.com/gorilla/websocket"
)

var chatUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var chatManager = ws.NewChatManager() // ✅ Ensure this exists in websocket package

// WebSocketChatHandler handles WebSocket connections for chat
func WebSocketChatHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := chatUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	// Register user for chat
	chatManager.HandleChatConnection(conn, userID)

	log.Printf("✅ User %d connected to WebSocket chat", userID)
}

// SendMessageHandler allows users to send messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	senderID := middlewares.GetUserIDFromSession(r)
	if senderID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var message models.ChatMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	message.SenderID = senderID

	db := config.GetDB()
	repo := repositories.NewChatRepository(db)

	err := repo.SaveMessage(&message)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetChatHistoryHandler retrieves chat history between two users
func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	otherUserID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || otherUserID == 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewChatRepository(db)

	messages, err := repo.GetMessages(userID, otherUserID)
	if err != nil {
		http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
