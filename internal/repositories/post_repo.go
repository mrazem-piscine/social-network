package repositories

import (
	"database/sql"
	"errors"
	"log"
	"social-network/internal/models"
)

// PostRepository handles post-related database operations
type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository creates a new instance of PostRepository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// GetUserPosts retrieves posts based on privacy settings
func (repo *PostRepository) GetUserPosts(userID int, viewerID int) ([]models.Post, error) {
	rows, err := repo.DB.Query(`
		SELECT id, user_id, content, image, privacy, created_at 
		FROM posts 
		WHERE user_id = ? 
		AND (privacy = 'public' OR 
		    (privacy = 'followers' AND user_id IN 
		        (SELECT following_id FROM followers WHERE follower_id = ? AND status = 'accepted')) 
		    OR user_id = ?) 
		ORDER BY created_at DESC`, userID, viewerID, viewerID)

	if err != nil {
		log.Println("Error fetching posts:", err)
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.Privacy, &post.CreatedAt)
		if err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// DeletePost deletes a post (only the creator can delete)
func (repo *PostRepository) DeletePost(postID, userID int) error {
	result, err := repo.DB.Exec(`
		DELETE FROM posts 
		WHERE id = ? AND user_id = ?`, postID, userID)

	if err != nil {
		log.Println("Error deleting post:", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("post not found or unauthorized")
	}

	return nil
}
