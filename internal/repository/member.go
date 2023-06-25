package repository

import (
	"chat-server/internal/domain/entity"
	"database/sql"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{
		db: db,
	}
}

func (r *MemberRepository) InsertMember(member *entity.Member) (*entity.Member, error) {
	query := "INSERT INTO members (room_id, user_id) VALUES ($1, $2)"
	_, err := r.db.Exec(query, member.RoomID, member.UserID)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *MemberRepository) SelectMembersByRoomID(roomID int) ([]*entity.Member, error) {
	query := "SELECT room_id, user_id FROM members WHERE room_id = $1"
	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*entity.Member
	for rows.Next() {
		var member entity.Member
		err := rows.Scan(&member.RoomID, &member.UserID)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

func (r *MemberRepository) UpdateMember(member *entity.Member) (*entity.Member, error) {
	query := "UPDATE members SET room_id = $1, user_id = $2 WHERE room_id = $3 AND user_id = $4"
	_, err := r.db.Exec(query, member.RoomID, member.UserID, member.RoomID, member.UserID)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *MemberRepository) DeleteMember(member *entity.Member) error {
	query := "DELETE FROM members WHERE room_id = $1 AND user_id = $2"
	_, err := r.db.Exec(query, member.RoomID, member.UserID)
	if err != nil {
		return err
	}
	return nil
}
