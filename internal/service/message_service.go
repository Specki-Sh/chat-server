package service

import (
	"fmt"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type MessageService struct {
	repo use_case.MessageStorage
}

func NewMessageService(repo use_case.MessageStorage) use_case.MessageUseCase {
	return &MessageService{
		repo: repo,
	}
}

func (m *MessageService) CreateMessage(req *entity.CreateMessageReq) (*entity.Message, error) {
	message := &entity.Message{
		SenderID: req.SenderID,
		RoomID:   req.RoomID,
		Content:  req.Content,
	}
	return m.repo.InsertMessage(message)
}

func (m *MessageService) GetMessageByID(id entity.ID) (*entity.Message, error) {
	return m.repo.SelectMessage(id)
}

func (m *MessageService) EditMessageContent(req *entity.EditMessageReq) (*entity.Message, error) {
	message, err := m.repo.SelectMessage(req.ID)
	if err != nil {
		return nil, fmt.Errorf("message service: %w", err)
	}
	message.Content = req.Content
	err = m.repo.UpdateMessage(message)
	if err != nil {
		return nil, fmt.Errorf("message service: %w", err)
	}
	return message, nil
}

func (m *MessageService) MarkReadMessageStatusByID(id entity.ID) error {
	message, err := m.repo.SelectMessage(id)
	if err != nil {
		return fmt.Errorf("message service: %w", err)
	}
	message.Status = "read"
	return m.repo.UpdateMessage(message)
}

func (m *MessageService) RemoveMessageByID(id entity.ID) error {
	if err := m.repo.SoftDeleteMessageByID(id); err != nil {
		return fmt.Errorf("message service: %w", err)
	}
	return nil
}

func (m *MessageService) GetMessagesPaginate(req *entity.GetMessagesPaginateReq) ([]*entity.Message, error) {
	messageBulk, err := m.repo.SelectMessagesPaginateReverse(req.RoomID, req.PerPage, req.Page)
	if err != nil {
		return nil, fmt.Errorf("message service: %w", err)
	}
	return messageBulk, nil
}

func (m *MessageService) IsMessageOwner(userID entity.ID, messageID entity.ID) (bool, error) {
	msg, err := m.repo.SelectMessage(messageID)
	if err != nil {
		return false, fmt.Errorf("message service: %w", err)
	}
	return msg.SenderID == userID, nil
}

func (m *MessageService) RemoveMessagesByRoomID(id entity.ID) error {
	if err := m.repo.SoftDeleteMessagesByRoomID(id); err != nil {
		return fmt.Errorf("message service: %w", err)
	}
	return nil
}
