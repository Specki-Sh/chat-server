package repository

import (
	"chat-server/internal/domain/entity"
	"database/sql"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) InsertMessage(message *entity.Message) (*entity.Message, error) {
	query := `INSERT INTO messages (sender_id, room_id, content) VALUES ($1, $2, $3) RETURNING id, timestamp`
	err := r.db.QueryRow(query, message.SenderID, message.RoomID, message.Content).Scan(&message.ID, &message.Timestamp)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *MessageRepository) SelectMessage(id int) (*entity.Message, error) {
	query := `SELECT id, sender_id, room_id, content, status, timestamp FROM messages WHERE id = $1`
	message := &entity.Message{}
	err := r.db.QueryRow(query, id).Scan(&message.ID, &message.SenderID, &message.RoomID,
		&message.Content, &message.Status, &message.Timestamp)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *MessageRepository) UpdateMessage(message *entity.Message) error {
	query := `UPDATE messages SET sender_id = $1, room_id = $2, content = $3, status = $4 WHERE id = $5`
	_, err := r.db.Exec(query, message.SenderID, message.RoomID, message.Content, message.Status, message.ID)
	return err
}

func (r *MessageRepository) DeleteMessage(id int) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MessageRepository) SelectMessagePaginate(perPage int, page int) ([]*entity.Message, error) {
	offset := perPage * (page - 1)
	query := `SELECT id, sender_id, room_id, content, status, timestamp FROM messages LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []*entity.Message
	for rows.Next() {
		message := &entity.Message{}
		err = rows.Scan(&message.ID, &message.SenderID, &message.RoomID, &message.Content, &message.Status, &message.Timestamp)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
