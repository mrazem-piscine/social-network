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

// CreateCommentHandler handles adding a comment to a post
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	// ✅ Get authenticated UserID from session
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// ✅ Ensure required fields are present
	if comment.Content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	commentRepo := repositories.NewCommentRepository(db)

	// ✅ Set UserID from session (prevents user spoofing)
	comment.UserID = userID

	// ✅ Save comment
	err := commentRepo.CreateComment(&comment)
	if err != nil {
		log.Println("Failed to add comment:", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	// ✅ Return JSON response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}
