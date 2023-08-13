package repository

import (
	"database/sql"
	"fmt"

	"chat-server/internal/domain/entity"
	dml "chat-server/pkg/db"
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
	err := m.db.QueryRow(query, message.SenderID, message.RoomID, message.Content).
		Scan(&message.ID, &message.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("MessageRepository.InsertMessage: %w", err)
	}
	return message, nil
}

func (m *MessageRepository) SelectMessage(id entity.ID) (*entity.Message, error) {
	query := dml.SelectMessageQuery
	message := &entity.Message{}
	err := m.db.QueryRow(query, id).Scan(&message.ID, &message.SenderID, &message.RoomID,
		&message.Content, &message.Status, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("MessageRepository.SelectMessage: %w", err)
	}
	return message, nil
}

func (m *MessageRepository) UpdateMessage(message *entity.Message) error {
	query := dml.UpdateMessageQuery
	_, err := m.db.Exec(
		query,
		message.SenderID,
		message.RoomID,
		message.Content,
		message.Status,
		message.ID,
	)
	if err != nil {
		return fmt.Errorf("MessageRepository.UpdateMessage: %w", err)
	}
	return nil
}

func (m *MessageRepository) SoftDeleteMessageByID(id entity.ID) error {
	query := dml.SoftDeleteMessageByIDQuery
	_, err := m.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("MessageRepository.SoftDeleteMessageByID: %w", err)
	}
	return nil
}

func (m *MessageRepository) SoftDeleteMessageBulkByRoomID(roomID entity.ID) error {
	query := dml.SoftDeleteMessageBulkByRoomIDQuery
	_, err := m.db.Exec(query, roomID)
	if err != nil {
		return fmt.Errorf("MessageRepository.SoftDeleteMessageBulkByRoomID: %w", err)
	}
	return nil
}

func (m *MessageRepository) SelectMessageBulkPaginate(
	roomID entity.ID,
	perPage uint,
	page uint,
) ([]entity.Message, error) {
	var messages []entity.Message
	offset := perPage * (page - 1)
	query := dml.SelectMessageBulkPaginateQuery
	rows, err := m.db.Query(query, roomID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("MessageRepository.SelectMessageBulkPaginate: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var message entity.Message
		err = rows.Scan(&message.ID, &message.SenderID, &message.RoomID, &message.Content,
			&message.Status, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("MessageRepository.SelectMessageBulkPaginate: %w", err)
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("MessageRepository.SelectMessageBulkPaginate: %w", err)
	}
	return messages, nil
}

func (m *MessageRepository) SelectMessageBulkPaginateReverse(
	roomID entity.ID,
	perPage uint,
	page uint,
) ([]entity.Message, error) {
	var messages []entity.Message
	offset := (page - 1) * perPage
	query := dml.SelectMessageBulkPaginateReverseQuery
	rows, err := m.db.Query(query, roomID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("MessageRepository.SelectMessageBulkPaginateReverse: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var message entity.Message
		err = rows.Scan(&message.ID, &message.SenderID, &message.RoomID, &message.Content,
			&message.Status, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("MessageRepository.SelectMessageBulkPaginateReverse: %w", err)
		}
		messages = append(messages, message)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("MessageRepository.SelectMessageBulkPaginateReverse: %w", err)
	}
	return messages, nil
}
