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

func (r *RoomService) CreateRoom(req *use_case.CreateRoomReq) (*use_case.CreateRoomRes, error) {
	room := entity.Room{
		OwnerID: req.OwnerID,
		Name:    req.Name,
	}
	newRoom, err := r.roomRepo.InsertRoom(&room)
	if err != nil {
		return nil, err
	}

	if _, err := r.AddMemberToRoom(newRoom.ID, newRoom.OwnerID); err != nil {
		return nil, err
	}

	res := use_case.CreateRoomRes{
		ID:      newRoom.ID,
		OwnerID: newRoom.OwnerID,
		Name:    newRoom.Name,
	}
	return &res, nil
}

func (r *RoomService) GetRoomInfoByID(id int) (*entity.Room, error) {
	return r.roomRepo.SelectRoomByID(id)
}

func (r *RoomService) EditRoomInfo(req *use_case.EditRoomReq) (*use_case.EditRoomRes, error) {
	room := entity.Room{
		ID:   req.ID,
		Name: req.Name,
	}
	err := r.roomRepo.UpdateRoom(&room)
	if err != nil {
		return nil, err
	}

	res := use_case.EditRoomRes{
		ID:   room.ID,
		Name: room.Name,
	}
	return &res, nil
}

func (r *RoomService) RemoveRoomByID(id int) error {
	return r.roomRepo.DeleteRoom(id)
}

func (r *RoomService) RoomExists(id int) (bool, error) {
	_, err := r.roomRepo.SelectRoomByID(id)
	if err != nil {
		if err == use_case.ErrRoomNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *RoomService) IsRoomOwner(roomID int, userID int) (bool, error) {
	room, err := r.roomRepo.SelectRoomByID(roomID)
	if err != nil {
		return false, err
	}
	return room.OwnerID == userID, nil
}

func (r *RoomService) HasRoomAccess(roomID int, userID int) (bool, error) {
	members, err := r.memberRepo.SelectMembersByRoomID(roomID)
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

func (r *RoomService) AddMemberToRoom(roomID int, userID int) (*entity.Member, error) {
	member := &entity.Member{RoomID: roomID, UserID: userID}
	return r.memberRepo.InsertMember(member)
}
