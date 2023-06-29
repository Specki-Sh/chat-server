package use_case

import "chat-server/internal/domain/entity"

type RoomUseCase interface {
	CreateRoom(req *entity.CreateRoomReq) (*entity.CreateRoomRes, error)
	GetRoomInfoByID(id int) (*entity.Room, error)
	EditRoomInfo(req *entity.EditRoomReq) (*entity.EditRoomRes, error)
	RemoveRoomByID(id int) error
	RoomExists(id int) (bool, error)
	IsRoomOwner(roomID int, userID int) (bool, error)
	HasRoomAccess(roomID int, userID int) (bool, error)
	AddMemberToRoom(roomID int, userID int) (*entity.Member, error)
}
