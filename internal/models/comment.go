package models

type Comment struct {
	ID       int    `json:"id"`
	PostID   int    `json:"post_id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}
