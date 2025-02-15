package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
)

// GroupChatRepository handles group chat-related database operations
type GroupChatRepository struct {
	DB *sql.DB
}

// NewGroupChatRepository creates a new instance of GroupChatRepository
func NewGroupChatRepository(db *sql.DB) *GroupChatRepository {
	return &GroupChatRepository{DB: db}
}

// SaveGroupMessage stores a group chat message in the database
func (repo *GroupChatRepository) SaveGroupMessage(message *models.GroupChatMessage) error {
	_, err := repo.DB.Exec(`
		INSERT INTO group_chat_messages (group_id, sender_id, content) 
		VALUES (?, ?, ?)`,
		message.GroupID, message.SenderID, message.Content)
	if err != nil {
		log.Println("Error saving group message:", err)
		return err
	}
	return nil
}

// GetGroupMessages fetches chat history for a group
func (repo *GroupChatRepository) GetGroupMessages(groupID int) ([]models.GroupChatMessage, error) {
	rows, err := repo.DB.Query(`
		SELECT id, group_id, sender_id, content, sent_at 
		FROM group_chat_messages 
		WHERE group_id = ? 
		ORDER BY sent_at ASC`, groupID)

	if err != nil {
		log.Println("Error fetching group messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.GroupChatMessage
	for rows.Next() {
		var msg models.GroupChatMessage
		if err := rows.Scan(&msg.ID, &msg.GroupID, &msg.SenderID, &msg.Content, &msg.SentAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
