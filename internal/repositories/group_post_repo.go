package repositories

import (
	"database/sql"
	"social-network/internal/models"
)

// GroupPostRepository handles group post-related operations
type GroupPostRepository struct {
	DB *sql.DB
}

// NewGroupPostRepository creates a new instance of GroupPostRepository
func NewGroupPostRepository(db *sql.DB) *GroupPostRepository {
	return &GroupPostRepository{DB: db}
}

// CreateGroupPost adds a new post to a group
func (repo *GroupPostRepository) CreateGroupPost(post *models.GroupPost) error {
	_, err := repo.DB.Exec(`
        INSERT INTO group_posts (group_id, user_id, content, image)
        VALUES (?, ?, ?, ?)`,
		post.GroupID, post.UserID, post.Content, post.Image)
	return err
}
