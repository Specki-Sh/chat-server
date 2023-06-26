package use_case

import (
	"chat-server/internal/domain/entity"
)

type CreateMessageReq struct {
	SenderID int    `json:"sender_id"`
	RoomID   int    `json:"room_id"`
	Content  string `json:"content"`
}

type EditMessageReq struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type GetMessagesPaginateReq struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
}

type MessageUseCase interface {
	CreateMessage(req *CreateMessageReq) (*entity.Message, error)
	GetMessageByID(id int) (*entity.Message, error)
	EditMessageContent(req *EditMessageReq) (*entity.Message, error)
	MarkReadMessageStatusByID(id int) error
	RemoveMessageByID(id int) error

	GetMessagesPaginate(req *GetMessagesPaginateReq) ([]*entity.Message, error)
}
