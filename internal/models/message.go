package models

import "time"

// DirectMessage represents a private message between users
type DirectMessage struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}
