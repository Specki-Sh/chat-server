package use_case

import "chat-server/internal/domain/entity"

type RoomUseCase interface {
	CreateRoom(req *entity.CreateRoomReq) (*entity.CreateRoomRes, error)
	GetRoomInfoByID(id entity.ID) (*entity.Room, error)
	EditRoomInfo(req *entity.EditRoomReq) (*entity.EditRoomRes, error)
	RemoveRoomByID(id entity.ID) error
	RoomExists(id entity.ID) (bool, error)
	IsRoomOwner(roomID entity.ID, userID entity.ID) (bool, error)
	HasRoomAccess(roomID entity.ID, userID entity.ID) (bool, error)
	AddMemberToRoom(roomID entity.ID, userID entity.ID) (*entity.Member, error)
}
