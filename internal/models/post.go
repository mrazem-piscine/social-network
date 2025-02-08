package models

import "time"

// Post represents a user post
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Image     *string   `json:"image"`
	Privacy   string    `json:"privacy"`
	CreatedAt time.Time `json:"created_at"`
}
