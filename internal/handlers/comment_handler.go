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

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	if comment.Content == "" || comment.PostID == 0 {
		http.Error(w, "Comment content and post ID are required", http.StatusBadRequest)
		return
	}

	comment.UserID = userID

	db := config.GetDB()
	commentRepo := repositories.NewCommentRepository(db)

	err := commentRepo.AddComment(&comment)
	if err != nil {
		log.Println("Error adding comment:", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment added successfully"})
}
func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
    userID := middlewares.GetUserIDFromSession(r)
    if userID == 0 {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
    if err != nil || commentID == 0 {
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }

    var requestBody struct {
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil || requestBody.Content == "" {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    db := config.GetDB()
    repo := repositories.NewCommentRepository(db)

    err = repo.EditComment(commentID, userID, requestBody.Content)
    if err != nil {
        log.Println("Error editing comment:", err)
        http.Error(w, "Failed to edit comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Comment updated successfully"})
}

func GetCommentsForPostHandler(w http.ResponseWriter, r *http.Request) {
    postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
    if err != nil || postID == 0 {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    db := config.GetDB()
    repo := repositories.NewCommentRepository(db)

    comments, err := repo.GetCommentsForPost(postID)
    if err != nil {
        log.Println("Error retrieving comments:", err)
        http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(comments)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
    userID := middlewares.GetUserIDFromSession(r)
    if userID == 0 {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
    if err != nil || commentID == 0 {
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }

    db := config.GetDB()
    repo := repositories.NewCommentRepository(db)

    err = repo.DeleteComment(commentID, userID)
    if err != nil {
        log.Println("Error deleting comment:", err)
        http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted successfully"})
}
