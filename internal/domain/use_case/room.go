package use_case

import "chat-server/internal/domain/entity"

type CreateRoomReq struct {
	OwnerId int    `json:"owner_id"`
	Name    string `json:"name"`
}

type CreateRoomRes struct {
	ID      int    `json:"id"`
	OwnerId int    `json:"owner_id"`
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
	GetRoomByID(id int) (*entity.Room, error)
	EditRoomInfo(req *EditRoomReq) (*EditRoomRes, error)
	RemoveRoomByID(id int) error
	RoomExists(id int) (bool, error)
}
