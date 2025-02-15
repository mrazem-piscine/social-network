package models

import "time"

// ChatMessage represents a message in the chat system
type ChatMessage struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID *int      `json:"receiver_id,omitempty"` // Nullable for group messages
	GroupID    *int      `json:"group_id,omitempty"`    // Nullable for direct messages
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}
