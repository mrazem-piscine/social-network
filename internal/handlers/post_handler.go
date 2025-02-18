package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/internal/config"
	"social-network/internal/middlewares"
	"social-network/internal/models"
	"social-network/internal/repositories"
	"strconv"
)

// CreatePostHandler allows users to create a post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil || post.Content == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	post.UserID = userID // Assign user ID to the post

	db := config.GetDB()
	repo := repositories.NewPostRepository(db)

	err := repo.CreatePost(&post)
	if err != nil {
		log.Println("Error creating post:", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

// GetUserPostsHandler retrieves posts based on user privacy settings
func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	viewingUserID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || viewingUserID == 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewPostRepository(db)
	posts, err := repo.GetUserPosts(viewingUserID, userID)
	if err != nil {
		log.Println("Error retrieving posts:", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil || postID == 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewPostRepository(db)

	err = repo.DeletePost(postID, userID)
	if err != nil {
		log.Println("Error deleting post:", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}

// EditPostHandler allows users to edit their own posts
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserIDFromSession(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil || postID == 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Content string  `json:"content"`
		Image   *string `json:"image"` // Nullable field
		Privacy string  `json:"privacy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db := config.GetDB()
	repo := repositories.NewPostRepository(db)

	err = repo.EditPost(postID, userID, requestBody.Content, requestBody.Image, requestBody.Privacy)
	if err != nil {
		log.Println("Error editing post:", err)
		http.Error(w, "Failed to edit post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
}
// GetAllPostsHandler retrieves all posts from all users
func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	db := config.GetDB()
	repo := repositories.NewPostRepository(db)

	posts, err := repo.GetAllPosts()
	if err != nil {
		log.Println("❌ Error retrieving all posts:", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
