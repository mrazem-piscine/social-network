package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
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
		log.Println("❌ Failed to upgrade WebSocket:", err)
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		return
	}

	// Register user for chat
	chatManager.HandleChatConnection(conn, userID)

	log.Printf("✅ User %d connected to WebSocket chat", userID)
}

// SendMessageHandler allows users to send messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var requestBody struct {
		ReceiverID int    `json:"receiver_id"`
		Content    string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("❌ Failed to decode request body:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if requestBody.ReceiverID == 0 || requestBody.Content == "" {
		log.Println("❌ Invalid message request: missing receiver or content")
		http.Error(w, "Invalid receiver ID or content", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewChatRepository(db) // ✅ FIXED function name

	err := repo.SaveMessage(userID, requestBody.ReceiverID, requestBody.Content)
	if err != nil {
		log.Println("❌ Failed to save message:", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	log.Printf("📩 Sending message: Sender %d -> Receiver %d: %s", userID, requestBody.ReceiverID, requestBody.Content)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
}

// GetChatHistoryHandler retrieves chat history between two users
func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverIDStr := r.URL.Query().Get("receiver_id")
	log.Printf("📌 Extracted receiverID from request: %s", receiverIDStr)

	receiverID, err := strconv.Atoi(receiverIDStr)
	if err != nil || receiverID == 0 {
		log.Println("❌ Invalid user ID:", receiverIDStr)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewChatRepository(db) // ✅ FIXED function name

	messages, err := repo.GetMessages(userID, receiverID) // ✅ FIXED function name
	if err != nil {
		log.Println("❌ Error retrieving chat history:", err)
		http.Error(w, "Failed to retrieve chat history", http.StatusInternalServerError)
		return
	}

	log.Printf("📜 Chat history retrieved between %d and %d", userID, receiverID)
	json.NewEncoder(w).Encode(messages)
}
