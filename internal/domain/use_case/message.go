package use_case

import (
	"chat-server/internal/domain/entity"
)

type MessageUseCase interface {
	CreateMessage(req *entity.CreateMessageReq) (*entity.Message, error)
	GetMessageByID(id int) (*entity.Message, error)
	EditMessageContent(req *entity.EditMessageReq) (*entity.Message, error)
	MarkReadMessageStatusByID(id int) error
	RemoveMessageByID(id int) error
	RemoveMessagesByRoomID(roomID int) error

	GetMessagesPaginate(req *entity.GetMessagesPaginateReq) ([]*entity.Message, error)
	IsMessageOwner(userID int, messageID int) (bool, error)
}
