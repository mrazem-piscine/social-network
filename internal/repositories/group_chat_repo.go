package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
	"social-network/internal/config"
)

// GroupChatRepository handles group chat-related database operations
type GroupChatRepository struct {
	DB *sql.DB
}

// NewGroupChatRepository creates a new instance of GroupChatRepository
func NewGroupChatRepository(db *sql.DB) *GroupChatRepository {
	return &GroupChatRepository{DB: db}
}

func (repo *GroupChatRepository) SaveGroupChatMessage(groupID, senderID int, content string) error {
		db := config.GetDB()
		_, err := db.Exec(`
			INSERT INTO group_chat_messages (group_id, sender_id, content)
			VALUES (?, ?, ?)`, groupID, senderID, content)
	
		if err != nil {
			log.Println("❌ Error saving message:", err)
		}
		return err
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
