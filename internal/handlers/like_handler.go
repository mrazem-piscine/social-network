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

// ToggleLikeHandler handles liking/unliking posts or comments
func ToggleLikeHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, _ := strconv.Atoi(r.URL.Query().Get("post_id"))
	commentID, _ := strconv.Atoi(r.URL.Query().Get("comment_id"))

	// Validate request
	if postID == 0 && commentID == 0 {
		http.Error(w, "Invalid request: Must provide post_id or comment_id", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewLikeRepository(db)

	liked, err := repo.ToggleLike(userID, postID, commentID)
	if err != nil {
		log.Println("Error toggling like:", err)
		http.Error(w, "Failed to like/unlike", http.StatusInternalServerError)
		return
	}

	status := "liked"
	if !liked {
		status = "unliked"
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

// GetLikeCountHandler retrieves the like count for a post or comment
func GetLikeCountHandler(w http.ResponseWriter, r *http.Request) {
	postID, _ := strconv.Atoi(r.URL.Query().Get("post_id"))
	commentID, _ := strconv.Atoi(r.URL.Query().Get("comment_id"))

	// Validate request
	if postID == 0 && commentID == 0 {
		http.Error(w, "Invalid request: Must provide post_id or comment_id", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewLikeRepository(db)

	likeCount, err := repo.GetLikeCount(postID, commentID)
	if err != nil {
		log.Println("Error retrieving like count:", err)
		http.Error(w, "Failed to get like count", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"like_count": likeCount})
}
