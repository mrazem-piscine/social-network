package models

import "time"

// GroupChatMessage represents a group chat message
type GroupChatMessage struct {
	ID       int       `json:"id"`
	GroupID  int       `json:"group_id"`
	SenderID int       `json:"sender_id"`
	Content  string    `json:"content"`
	SentAt   time.Time `json:"sent_at"`
}
