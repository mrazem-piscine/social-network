package models

// GroupPost represents a post inside a group
type GroupPost struct {
	ID        int     `json:"id"`
	GroupID   int     `json:"group_id"`
	UserID    int     `json:"user_id"`
	Content   string  `json:"content"`
	Image     *string `json:"image,omitempty"`
	CreatedAt string  `json:"created_at"`
}
