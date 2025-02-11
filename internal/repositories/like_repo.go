package repositories

import (
	"database/sql"
	"log"
)

// LikeRepository handles database operations for likes
type LikeRepository struct {
	DB *sql.DB
}

// NewLikeRepository creates a new instance of LikeRepository
func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{DB: db}
}

// ToggleLike adds or removes a like for a post or comment
func (repo *LikeRepository) ToggleLike(userID, postID, commentID int) (bool, error) {
	// Check if the like already exists
	var exists bool
	err := repo.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id IS ? AND comment_id IS ?)`,
		userID, postID, commentID).Scan(&exists)

	if err != nil {
		log.Println("Error checking like existence:", err)
		return false, err
	}

	if exists {
		// Remove the like
		_, err := repo.DB.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id IS ? AND comment_id IS ?`,
			userID, postID, commentID)
		if err != nil {
			log.Println("Error removing like:", err)
			return false, err
		}
		return false, nil // Indicates unlike
	}

	// Add the like
	_, err = repo.DB.Exec(`INSERT INTO likes (user_id, post_id, comment_id) VALUES (?, ?, ?)`,
		userID, postID, commentID)
	if err != nil {
		log.Println("Error inserting like:", err)
		return false, err
	}

	return true, nil // Indicates like added
}

// GetLikeCount retrieves the number of likes for a post or comment
func (repo *LikeRepository) GetLikeCount(postID, commentID int) (int, error) {
	var count int
	err := repo.DB.QueryRow(`
		SELECT COUNT(*) FROM likes WHERE post_id IS ? AND comment_id IS ?`,
		postID, commentID).Scan(&count)

	if err != nil {
		log.Println("Error counting likes:", err)
		return 0, err
	}

	return count, nil
}
