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
	msg, err := m.repo.InsertMessage(message)
	if err != nil {
		return nil, fmt.Errorf("MesssageService.CreateMessage: %w", err)
	}
	return msg, nil
}

func (m *MessageService) GetMessageByID(id entity.ID) (*entity.Message, error) {
	msg, err := m.repo.SelectMessage(id)
	if err != nil {
		return nil, fmt.Errorf("MesssageService.GetMessageByID: %w", err)
	}
	return msg, nil
}

func (m *MessageService) EditMessageContent(req *entity.EditMessageReq) (*entity.Message, error) {
	message, err := m.repo.SelectMessage(req.ID)
	if err != nil {
		return nil, fmt.Errorf("MesssageService.EditMessageContent: %w", err)
	}
	message.Content = req.Content
	err = m.repo.UpdateMessage(message)
	if err != nil {
		return nil, fmt.Errorf("MesssageService.EditMessageContent: %w", err)
	}
	return message, nil
}

func (m *MessageService) MarkReadMessageStatusByID(id entity.ID) error {
	message, err := m.repo.SelectMessage(id)
	if err != nil {
		return fmt.Errorf("MesssageService.MarkReadMessageStatusByID: %w", err)
	}
	message.Status = "read"
	if err := m.repo.UpdateMessage(message); err != nil {
		return fmt.Errorf("MesssageService.MarkReadMessageStatusByID: err")
	}
	return nil
}

func (m *MessageService) RemoveMessageByID(id entity.ID) error {
	if err := m.repo.SoftDeleteMessageByID(id); err != nil {
		return fmt.Errorf("MesssageService.RemoveMessageByID: %w", err)
	}
	return nil
}

func (m *MessageService) GetMessageBulkPaginate(
	req *entity.GetMessageBulkPaginateReq,
) ([]entity.Message, error) {
	messageBulk, err := m.repo.SelectMessageBulkPaginateReverse(req.RoomID, req.PerPage, req.Page)
	if err != nil {
		return nil, fmt.Errorf("MesssageService.GetMessageBulkPaginate: %w", err)
	}
	return messageBulk, nil
}

func (m *MessageService) IsMessageOwner(userID entity.ID, messageID entity.ID) (bool, error) {
	msg, err := m.repo.SelectMessage(messageID)
	if err != nil {
		return false, fmt.Errorf("MesssageService.IsMessageOwner: %w", err)
	}
	return msg.SenderID == userID, nil
}

func (m *MessageService) RemoveMessageBulkByRoomID(id entity.ID) error {
	if err := m.repo.SoftDeleteMessageBulkByRoomID(id); err != nil {
		return fmt.Errorf("MesssageService.RemoveMessageBulkByRoomID: %w", err)
	}
	return nil
}
