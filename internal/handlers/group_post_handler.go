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

// CreateGroupPostHandler handles posting inside a group
func CreateGroupPostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var post models.GroupPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if post.Content == "" || post.GroupID == 0 {
		http.Error(w, "Group ID and content are required", http.StatusBadRequest)
		return
	}

	post.UserID = userID

	db := config.GetDB()
	repo := repositories.NewGroupPostRepository(db)

	err := repo.CreateGroupPost(&post)
	if err != nil {
		log.Println("Error creating group post:", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group post created successfully"})
}
