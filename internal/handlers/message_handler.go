package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/repositories"
	"strconv"
)

// GetMessagesHandler retrieves chat history between two users
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	receiverID, err := strconv.Atoi(r.URL.Query().Get("receiver_id"))
	if err != nil || receiverID == 0 {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewMessageRepository(db)

	messages, err := repo.GetMessages(userID, receiverID)
	if err != nil {
		log.Println("Error retrieving messages:", err)
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
