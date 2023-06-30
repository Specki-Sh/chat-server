package repository

import (
	"chat-server/internal/domain/entity"
	dml "chat-server/pkg/db"
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

func (m *MessageRepository) InsertMessage(message *entity.Message) (*entity.Message, error) {
	query := dml.InsertMessageQuery
	err := m.db.QueryRow(query, message.SenderID, message.RoomID, message.Content).Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *MessageRepository) SelectMessage(id entity.ID) (*entity.Message, error) {
	query := dml.SelectMessageQuery
	message := &entity.Message{}
	err := m.db.QueryRow(query, id).Scan(&message.ID, &message.SenderID, &message.RoomID,
		&message.Content, &message.Status, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *MessageRepository) UpdateMessage(message *entity.Message) error {
	query := dml.UpdateMessageQuery
	_, err := m.db.Exec(query, message.SenderID, message.RoomID, message.Content, message.Status, message.ID)
	return err
}

func (m *MessageRepository) SoftDeleteMessageByID(id entity.ID) error {
	query := dml.SoftDeleteMessageByIDQuery
	_, err := m.db.Exec(query, id)
	return err
}

func (m *MessageRepository) SoftDeleteMessagesByRoomID(roomID entity.ID) error {
	query := dml.SoftDeleteMessagesByRoomIDQuery
	_, err := m.db.Exec(query, roomID)
	return err
}

func (m *MessageRepository) SelectMessagePaginate(roomID entity.ID, perPage uint, page uint) ([]*entity.Message, error) {
	var messages []*entity.Message
	offset := perPage * (page - 1)
	query := dml.SelectMessagePaginateQuery
	rows, err := m.db.Query(query, roomID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		message := &entity.Message{}
		err = rows.Scan(&message.ID, &message.SenderID, &message.RoomID, &message.Content, &message.Status, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
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

func (m *MessageRepository) SelectMessagesPaginateReverse(roomID entity.ID, perPage uint, page uint) ([]*entity.Message, error) {
	var messages []*entity.Message
	offset := (page - 1) * perPage
	query := dml.SelectMessagesPaginateReverseQuery
	rows, err := m.db.Query(query, roomID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var message entity.Message
		err = rows.Scan(&message.ID, &message.SenderID, &message.RoomID, &message.Content, &message.Status, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
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
