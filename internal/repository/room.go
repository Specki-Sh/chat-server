package repository

import (
	"database/sql"
	"fmt"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	dml "chat-server/pkg/db"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r *RoomRepository) InsertRoom(room *entity.Room) (*entity.Room, error) {
	if room.OwnerID == 0 || room.Name == "" {
		return nil, fmt.Errorf("RoomRepository.InsertRoom: %w", use_case.ErrRoomInvalid)
	}
	query := dml.InsertRoomQuery
	err := r.db.QueryRow(query, room.OwnerID, room.Name).Scan(&room.ID)
	if err != nil {
		return nil, fmt.Errorf("RoomRepository.InsertRoom: %w", err)
	}
	return room, nil
}

func (r *RoomRepository) SelectRoomByID(id entity.ID) (*entity.Room, error) {
	query := dml.SelectRoomByIDQuery
	var room entity.Room
	err := r.db.QueryRow(query, id).Scan(&room.ID, &room.OwnerID, &room.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("RoomRepository.SelectRoomByID: %w", use_case.ErrRoomNotFound)
		}
		return nil, fmt.Errorf("RoomRepository.SelectRoomByID: %w", err)
	}
	return &room, nil
}

func (r *RoomRepository) UpdateRoom(room *entity.Room) error {
	if room.ID == 0 || room.Name == "" {
		return fmt.Errorf("RoomRepository.UpdateRoom: %w", use_case.ErrRoomInvalid)
	}
	query := dml.UpdateRoomQuery
	res, err := r.db.Exec(query, room.Name, room.ID)
	if err != nil {
		return fmt.Errorf("RoomRepository.UpdateRoom: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RoomRepository.UpdateRoom: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("RoomRepository.UpdateRoom: %w", use_case.ErrRoomNotFound)
	}
	return nil
}

func (r *RoomRepository) DeleteRoom(id entity.ID) error {
	query := dml.DeleteRoomQuery
	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("RoomRepository.DeleteRoom: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RoomRepository.DeleteRoom: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("RoomRepository.DeleteRoom: %w", use_case.ErrRoomNotFound)
	}
	return nil
}
