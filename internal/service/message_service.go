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

func (s *MessageService) CreateMessage(req *use_case.CreateMessageReq) (*entity.Message, error) {
	message := &entity.Message{
		SenderID: req.SenderID,
		RoomID:   req.RoomID,
		Content:  req.Content,
	}
	return s.repo.InsertMessage(message)
}

func (s *MessageService) GetMessageByID(id int) (*entity.Message, error) {
	return s.repo.SelectMessage(id)
}

func (s *MessageService) EditMessageContent(req *use_case.EditMessageReq) (*entity.Message, error) {
	message, err := s.repo.SelectMessage(req.ID)
	if err != nil {
		return nil, err
	}
	message.Content = req.Content
	err = s.repo.UpdateMessage(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (s *MessageService) MarkReadMessageStatusByID(id int) error {
	message, err := s.repo.SelectMessage(id)
	if err != nil {
		return err
	}
	message.Status = "read"
	return s.repo.UpdateMessage(message)
}

func (s *MessageService) RemoveMessageByID(id int) error {
	return s.repo.SoftDeleteMessage(id)
}

func (s *MessageService) GetMessagesPaginate(req *use_case.GetMessagesPaginateReq) ([]*entity.Message, error) {
	return s.repo.SelectMessagesPaginateReverse(req.RoomID, req.PerPage, req.Page)
}

func (s *MessageService) IsMessageOwner(userID int, messageID int) (bool, error) {
	msg, err := s.repo.SelectMessage(messageID)
	if err != nil {
		return false, err
	}
	return msg.SenderID == userID, nil
}
