package repository

import (
	"chat-server/internal/domain/use_case"
	"database/sql"

	"chat-server/internal/domain/entity"
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
		return nil, use_case.ErrRoomInvalid
	}
	query := "INSERT INTO rooms (owner_id, name) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, room.OwnerID, room.Name).Scan(&room.ID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomRepository) SelectRoomByID(id entity.ID) (*entity.Room, error) {
	query := "SELECT id, owner_id, name FROM rooms WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var room entity.Room
	err := row.Scan(&room.ID, &room.OwnerID, &room.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, use_case.ErrRoomNotFound
		}
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) UpdateRoom(room *entity.Room) error {
	if room.ID == 0 || room.Name == "" {
		return use_case.ErrRoomInvalid
	}
	query := "UPDATE rooms SET name = $1 WHERE id = $2"
	res, err := r.db.Exec(query, room.Name, room.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return use_case.ErrRoomNotFound
	}
	return nil
}

func (r *RoomRepository) DeleteRoom(id entity.ID) error {
	query := "DELETE FROM rooms WHERE id = $1"
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return use_case.ErrRoomNotFound
	}
	return nil
}
