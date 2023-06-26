package service

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type RoomService struct {
	roomRepo   use_case.RoomRepository
	memberRepo use_case.MemberRepository
}

func NewRoomService(roomRepo use_case.RoomRepository, memberRepo use_case.MemberRepository) *RoomService {
	return &RoomService{
		roomRepo:   roomRepo,
		memberRepo: memberRepo,
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
	return s.roomRepo.SelectRoomByID(id)
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

func (s *RoomService) RoomExists(id int) (bool, error) {
	_, err := s.roomRepo.SelectRoomByID(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *RoomService) IsRoomOwner(roomID int, userID int) (bool, error) {
	room, err := s.roomRepo.SelectRoomByID(roomID)
	if err != nil {
		return false, err
	}
	return room.OwnerID == userID, nil
}

func (s *RoomService) HasRoomAccess(roomID int, userID int) (bool, error) {
	members, err := s.memberRepo.SelectMembersByRoomID(roomID)
	if err != nil {
		return false, err
	}
	for _, member := range members {
		if member.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (s *RoomService) AddMemberToRoom(roomID int, userID int) (*entity.Member, error) {
	member := &entity.Member{RoomID: roomID, UserID: userID}
	return s.memberRepo.InsertMember(member)
}
