package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
)

// ChatRepository handles chat-related database operations
type ChatRepository struct {
	DB *sql.DB
}

// NewChatRepository creates a new instance of ChatRepository
func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{DB: db}
}

// SaveMessage stores a chat message in the database
func (repo *ChatRepository) SaveMessage(message *models.ChatMessage) error {
	_, err := repo.DB.Exec(`
		INSERT INTO chat_messages (sender_id, receiver_id, group_id, content, sent_at) 
		VALUES (?, ?, ?, ?, ?)`,
		message.SenderID, message.ReceiverID, message.GroupID, message.Content, message.SentAt)
	if err != nil {
		log.Println("Error saving message:", err)
		return err
	}
	return nil
}

// GetMessages fetches chat history between two users
func (repo *ChatRepository) GetMessages(user1, user2 int) ([]models.ChatMessage, error) {
	rows, err := repo.DB.Query(`
		SELECT id, sender_id, receiver_id, group_id, content, sent_at 
		FROM chat_messages 
		WHERE (sender_id = ? AND receiver_id = ?) 
		OR (sender_id = ? AND receiver_id = ?) 
		ORDER BY sent_at ASC`, user1, user2, user2, user1)

	if err != nil {
		log.Println("Error fetching messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.GroupID, &msg.Content, &msg.SentAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
