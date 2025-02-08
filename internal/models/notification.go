package models

type Notification struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}
