package service

import (
	"fmt"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
)

type roomService struct {
	roomRepo   use_case.RoomStorage
	memberRepo use_case.MemberStorage
}

func NewRoomService(roomRepo use_case.RoomStorage, memberRepo use_case.MemberStorage) use_case.RoomUseCase {
	return &roomService{
		roomRepo:   roomRepo,
		memberRepo: memberRepo,
	}
}

func (r *roomService) CreateRoom(req *entity.CreateRoomReq) (*entity.CreateRoomRes, error) {
	room := entity.Room{
		OwnerID: req.OwnerID,
		Name:    req.Name,
	}
	newRoom, err := r.roomRepo.InsertRoom(&room)
	if err != nil {
		return nil, fmt.Errorf("roomService.CreateRoom: %w", err)
	}

	if _, err := r.AddMemberToRoom(newRoom.ID, newRoom.OwnerID); err != nil {
		return nil, fmt.Errorf("roomService.CreateRoom: %w", err)
	}

	res := entity.CreateRoomRes{
		ID:      newRoom.ID,
		OwnerID: newRoom.OwnerID,
		Name:    newRoom.Name,
	}
	return &res, nil
}

func (r *roomService) GetRoomInfoByID(id entity.ID) (*entity.Room, error) {
	room, err := r.roomRepo.SelectRoomByID(id)
	if err != nil {
		return nil, fmt.Errorf("roomService.GetRoomInfoByID: %w", err)
	}
	return room, nil
}

func (r *roomService) EditRoomInfo(req *entity.EditRoomReq) (*entity.EditRoomRes, error) {
	room := entity.Room{
		ID:   req.ID,
		Name: req.Name,
	}
	err := r.roomRepo.UpdateRoom(&room)
	if err != nil {
		return nil, fmt.Errorf("roomService.EditRoomInfo: %w", err)
	}

	res := entity.EditRoomRes{
		ID:   room.ID,
		Name: room.Name,
	}
	return &res, nil
}

func (r *roomService) RemoveRoomByID(id entity.ID) error {
	if err := r.roomRepo.DeleteRoom(id); err != nil {
		return fmt.Errorf("roomService.RemoveRoomByID: %w", err)
	}
	return nil
}

func (r *roomService) RoomExists(id entity.ID) (bool, error) {
	_, err := r.roomRepo.SelectRoomByID(id)
	if err != nil {
		if err == use_case.ErrRoomNotFound {
			return false, nil
		}
		return false, fmt.Errorf("roomService.RoomExists: %w", err)
	}
	return true, nil
}

func (r *roomService) IsRoomOwner(roomID entity.ID, userID entity.ID) (bool, error) {
	room, err := r.roomRepo.SelectRoomByID(roomID)
	if err != nil {
		return false, fmt.Errorf("roomService.IsRoomOwner: %w", err)
	}
	return room.OwnerID == userID, nil
}

func (r *roomService) HasRoomAccess(roomID entity.ID, userID entity.ID) (bool, error) {
	members, err := r.memberRepo.SelectMemberBulkByRoomID(roomID)
	if err != nil {
		return false, fmt.Errorf("roomService.HasRoomAccess: %w", err)
	}
	for _, member := range members {
		if member.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (r *roomService) AddMemberToRoom(roomID entity.ID, userID entity.ID) (*entity.Member, error) {
	member := &entity.Member{RoomID: roomID, UserID: userID}
	m, err := r.memberRepo.InsertMember(member)
	if err != nil {
		return nil, fmt.Errorf("roomService.AddMemberToRoom: %w", err)
	}
	return m, nil
}
