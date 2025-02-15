package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	"social-network/internal/repositories"
)

// CreateGroupEventHandler handles event creation inside a group
func CreateGroupEventHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var event models.GroupEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if event.Title == "" || event.Description == "" || event.GroupID == 0 || event.EventDate == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	event.CreatorID = userID

	db := config.GetDB()
	repo := repositories.NewGroupEventRepository(db)

	err := repo.CreateGroupEvent(&event)
	if err != nil {
		log.Println("Error creating group event:", err)
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group event created successfully"})
}
