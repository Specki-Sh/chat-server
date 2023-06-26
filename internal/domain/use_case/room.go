package use_case

import "chat-server/internal/domain/entity"

type CreateRoomReq struct {
	OwnerID int    `json:"owner_id"`
	Name    string `json:"name"`
}

type CreateRoomRes struct {
	ID      int    `json:"id"`
	OwnerID int    `json:"owner_id"`
	Name    string `json:"name"`
}

type EditRoomReq struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EditRoomRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoomUseCase interface {
	CreateRoom(req *CreateRoomReq) (*CreateRoomRes, error)
	GetRoomInfoByID(id int) (*entity.Room, error)
	EditRoomInfo(req *EditRoomReq) (*EditRoomRes, error)
	RemoveRoomByID(id int) error
	RoomExists(id int) (bool, error)
	IsRoomOwner(roomID int, userID int) (bool, error)
	HasRoomAccess(roomID int, userID int) (bool, error)
	AddMemberToRoom(roomID int, userID int) (*entity.Member, error)
}
