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
	query := `INSERT INTO messages (sender_id, room_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(query, message.SenderID, message.RoomID, message.Content).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *MessageRepository) SelectMessage(id int) (*entity.Message, error) {
	query := `SELECT id, sender_id, room_id, content, status, created_at FROM messages WHERE id = $1 AND is_active = true`
	message := &entity.Message{}
	err := r.db.QueryRow(query, id).Scan(&message.ID, &message.SenderID, &message.RoomID,
		&message.Content, &message.Status, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (r *MessageRepository) UpdateMessage(message *entity.Message) error {
	query := `UPDATE messages SET sender_id = $1, room_id = $2, content = $3, status = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	_, err := r.db.Exec(query, message.SenderID, message.RoomID, message.Content, message.Status, message.ID)
	return err
}

func (r *MessageRepository) SoftDeleteMessageByID(id int) error {
	query := `UPDATE messages SET is_active = false, deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *MessageRepository) SoftDeleteMessagesByRoomID(roomID int) error {
	query := `UPDATE messages SET is_active = false, deleted_at = CURRENT_TIMESTAMP WHERE room_id = $1`
	_, err := r.db.Exec(query, roomID)
	return err
}

func (r *MessageRepository) SelectMessagePaginate(roomID int, perPage int, page int) ([]*entity.Message, error) {
	var messages []*entity.Message
	offset := perPage * (page - 1)
	query := `SELECT id, sender_id, room_id, content, status, created_at FROM messages WHERE is_active = true AND room_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, roomID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		message := &entity.Message{}
		err = rows.Scan(&message.ID, &message.SenderID, &message.RoomID, &message.Content, &message.Status, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageRepository) SelectMessagesPaginateReverse(roomID int, perPage int, page int) ([]*entity.Message, error) {
	var messages []*entity.Message
	offset := (page - 1) * perPage
	query := `SELECT id,sender_id,room_id,content,status,created_at FROM messages WHERE is_active = true AND room_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, roomID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message entity.Message
		err = rows.Scan(&message.ID, &message.SenderID, &message.RoomID, &message.Content, &message.Status, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return messages, nil
}
