package repository

import (
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
	query := "INSERT INTO rooms (owner_id, name) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRow(query, room.OwnerID, room.Name).Scan(&room.Id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomRepository) SelectRoom(id int) (*entity.Room, error) {
	query := "SELECT id, owner_id, name FROM rooms WHERE id = $1"
	row := r.db.QueryRow(query, id)

	var room entity.Room
	err := row.Scan(&room.Id, &room.OwnerID, &room.Name)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) UpdateRoom(room *entity.Room) error {
	query := "UPDATE rooms SET name = $1 WHERE id = $2"
	_, err := r.db.Exec(query, room.Name, room.Id)
	return err
}

func (r *RoomRepository) DeleteRoom(id int) error {
	query := "DELETE FROM rooms WHERE id = $1"
	_, err := r.db.Exec(query, id)
	return err
}
