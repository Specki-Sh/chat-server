package service

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type RoomService struct {
	roomRepo use_case.RoomRepository
}

func NewRoomService(roomRepo use_case.RoomRepository) *RoomService {
	return &RoomService{
		roomRepo: roomRepo,
	}
}

func (s *RoomService) CreateRoom(req *use_case.CreateRoomReq) (*use_case.CreateRoomRes, error) {
	room := entity.Room{
		OwnerID: req.OwnerId,
		Name:    req.Name,
	}
	newRoom, err := s.roomRepo.InsertRoom(&room)
	if err != nil {
		return nil, err
	}

	res := use_case.CreateRoomRes{
		ID:      newRoom.Id,
		OwnerId: newRoom.OwnerID,
		Name:    newRoom.Name,
	}
	return &res, nil
}

func (s *RoomService) GetRoomByID(id int) (*entity.Room, error) {
	return s.roomRepo.SelectRoom(id)
}

func (s *RoomService) EditRoomInfo(req *use_case.EditRoomReq) (*use_case.EditRoomRes, error) {
	room := entity.Room{
		Id:   req.ID,
		Name: req.Name,
	}
	err := s.roomRepo.UpdateRoom(&room)
	if err != nil {
		return nil, err
	}

	res := use_case.EditRoomRes{
		ID:   room.Id,
		Name: room.Name,
	}
	return &res, nil
}

func (s *RoomService) RemoveRoomByID(id int) error {
	return s.roomRepo.DeleteRoom(id)
}
