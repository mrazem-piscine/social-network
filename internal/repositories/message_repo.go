package repositories

import (
	"database/sql"
	"log"
	"social-network/internal/models"
)

// MessageRepository handles message operations
type MessageRepository struct {
	DB *sql.DB
}

// NewMessageRepository creates a new instance
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

// SendMessage inserts a new message into the database
func (repo *MessageRepository) SendMessage(msg *models.DirectMessage) error {
	_, err := repo.DB.Exec(`
		INSERT INTO direct_messages (sender_id, receiver_id, content)
		VALUES (?, ?, ?)`,
		msg.SenderID, msg.ReceiverID, msg.Content,
	)
	return err
}

// GetMessages retrieves messages between two users
func (repo *MessageRepository) GetMessages(senderID, receiverID int) ([]models.DirectMessage, error) {
	rows, err := repo.DB.Query(`
		SELECT id, sender_id, receiver_id, content, sent_at 
		FROM direct_messages 
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY sent_at ASC`,
		senderID, receiverID, receiverID, senderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.DirectMessage
	for rows.Next() {
		var msg models.DirectMessage
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.SentAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
func (repo *MessageRepository) SaveMessage(senderID, receiverID int, content string) error {
	_, err := repo.DB.Exec(`
        INSERT INTO direct_messages (sender_id, receiver_id, content)
        VALUES (?, ?, ?)`,
		senderID, receiverID, content)

	if err != nil {
		log.Println("❌ Error inserting message into database:", err)
	} else {
		log.Printf("✅ Message saved: %d -> %d: %s", senderID, receiverID, content)
	}
	return err
}
func (repo *MessageRepository) GetChatHistory(userID, receiverID int) ([]models.DirectMessage, error) {
	rows, err := repo.DB.Query(`
        SELECT id, sender_id, receiver_id, content, sent_at
        FROM direct_messages
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY sent_at ASC`,
		userID, receiverID, receiverID, userID)

	if err != nil {
		log.Println("❌ Database error retrieving chat history:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.DirectMessage
	for rows.Next() {
		var message models.DirectMessage
		err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.SentAt)
		if err != nil {
			log.Println("❌ Error scanning row:", err)
			return nil, err
		}
		messages = append(messages, message)
	}
	log.Printf("📌 Fetching chat history for User %d ↔ User %d", userID, receiverID)

	return messages, nil
}
