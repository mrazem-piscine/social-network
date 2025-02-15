package models

import "time"

// GroupMember represents a user in a group
type GroupMember struct {
	UserID   int       `json:"user_id"`
	Nickname string    `json:"nickname"`
	Status   string    `json:"status"`    // Fix: Add "status" field
	JoinedAt time.Time `json:"joined_at"` // Fix: Add "joined_at" field
}

type Group struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatorID   int       `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}
