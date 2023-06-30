package use_case

import (
	"chat-server/internal/domain/entity"
)

type MessageUseCase interface {
	CreateMessage(req *entity.CreateMessageReq) (*entity.Message, error)
	GetMessageByID(id entity.ID) (*entity.Message, error)
	EditMessageContent(req *entity.EditMessageReq) (*entity.Message, error)
	MarkReadMessageStatusByID(id entity.ID) error
	RemoveMessageByID(id entity.ID) error
	RemoveMessagesByRoomID(roomID entity.ID) error

	GetMessagesPaginate(req *entity.GetMessagesPaginateReq) ([]*entity.Message, error)
	IsMessageOwner(userID entity.ID, messageID entity.ID) (bool, error)
}
