package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
	"time"
)

// PostRepository handles post-related database operations
type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository creates a new instance of PostRepository
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}
func (repo *PostRepository) CreatePost(post *models.Post) error {
	_, err := repo.DB.Exec(`
    INSERT INTO posts (user_id, content, image, privacy, created_at)
    VALUES (?, ?, ?, ?, ?)`,
		post.UserID, post.Content, post.Image, post.Privacy, time.Now(),
	)

	return err
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
	_, err := repo.DB.Exec(`DELETE FROM posts WHERE id = ? AND user_id = ?`, postID, userID)
	return err
}

// EditPost updates a post's content
func (repo *PostRepository) EditPost(postID, userID int, content string, image *string, privacy string) error {
	_, err := repo.DB.Exec(`
        UPDATE posts
        SET content = ?, image = ?, privacy = ?
        WHERE id = ? AND user_id = ?`,
		content, image, privacy, postID, userID,
	)
	return err
}
