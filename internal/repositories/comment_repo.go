package repositories

import (
	"database/sql"
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

func (repo *CommentRepository) AddComment(comment *models.Comment) error {
	_, err := repo.DB.Exec(`
        INSERT INTO comments (post_id, user_id, content, created_at) 
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
		comment.PostID, comment.UserID, comment.Content,
	)
	return err
}

func (repo *CommentRepository) DeleteComment(commentID, userID int) error {
	_, err := repo.DB.Exec(`DELETE FROM comments WHERE id = ? AND user_id = ?`, commentID, userID)
	return err
}

func (repo *CommentRepository) GetCommentsForPost(postID int) ([]models.Comment, error) {
	var comments []models.Comment

	rows, err := repo.DB.Query(`
        SELECT id, post_id, user_id, content, created_at
        FROM comments
        WHERE post_id = ?
        ORDER BY created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
func (repo *CommentRepository) EditComment(commentID, userID int, content string) error {
	_, err := repo.DB.Exec(`
        UPDATE comments
        SET content = ?
        WHERE id = ? AND user_id = ?`,
		content, commentID, userID,
	)
	return err
}
