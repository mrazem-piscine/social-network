package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
)

// CommentRepository handles comment-related database operations
type CommentRepository struct {
	DB *sql.DB
}

// NewCommentRepository creates a new instance of CommentRepository
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// CreateComment inserts a comment into the database
func (repo *CommentRepository) CreateComment(comment *models.Comment) error {
	_, err := repo.DB.Exec(`
		INSERT INTO comments (post_id, user_id, username, content) 
		VALUES (?, ?, ?, ?)`,

		comment.PostID, comment.UserID, comment.Username, comment.Content,
	)
	if err != nil {
		log.Println("Error inserting comment:", err)
	}
	return err
}
