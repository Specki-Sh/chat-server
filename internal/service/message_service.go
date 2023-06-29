package service

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type MessageService struct {
	repo use_case.MessageRepository
}

func NewMessageService(repo use_case.MessageRepository) *MessageService {
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

func (m *MessageService) GetMessageByID(id int) (*entity.Message, error) {
	return m.repo.SelectMessage(id)
}

func (m *MessageService) EditMessageContent(req *entity.EditMessageReq) (*entity.Message, error) {
	message, err := m.repo.SelectMessage(req.ID)
	if err != nil {
		return nil, err
	}
	message.Content = req.Content
	err = m.repo.UpdateMessage(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *MessageService) MarkReadMessageStatusByID(id int) error {
	message, err := m.repo.SelectMessage(id)
	if err != nil {
		return err
	}
	message.Status = "read"
	return m.repo.UpdateMessage(message)
}

func (m *MessageService) RemoveMessageByID(id int) error {
	return m.repo.SoftDeleteMessageByID(id)
}

func (m *MessageService) GetMessagesPaginate(req *entity.GetMessagesPaginateReq) ([]*entity.Message, error) {
	return m.repo.SelectMessagesPaginateReverse(req.RoomID, req.PerPage, req.Page)
}

func (m *MessageService) IsMessageOwner(userID int, messageID int) (bool, error) {
	msg, err := m.repo.SelectMessage(messageID)
	if err != nil {
		return false, err
	}
	return msg.SenderID == userID, nil
}

func (m *MessageService) RemoveMessagesByRoomID(id int) error {
	return m.repo.SoftDeleteMessagesByRoomID(id)
}
