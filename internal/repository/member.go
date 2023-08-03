package repository

import (
	"database/sql"

	"chat-server/internal/domain/entity"
	dml "chat-server/pkg/db"
)

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{
		db: db,
	}
}

func (m *MemberRepository) InsertMember(member *entity.Member) (*entity.Member, error) {
	query := dml.InsertMemberQuery
	_, err := m.db.Exec(query, member.RoomID, member.UserID)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MemberRepository) SelectMemberBulkByRoomID(roomID entity.ID) ([]entity.Member, error) {
	query := dml.SelectMemberBulkByRoomIDQuery
	rows, err := m.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []entity.Member
	for rows.Next() {
		var member entity.Member
		err := rows.Scan(&member.RoomID, &member.UserID)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (m *MemberRepository) UpdateMember(member *entity.Member) (*entity.Member, error) {
	query := dml.UpdateMemberQuery
	_, err := m.db.Exec(query, member.RoomID, member.UserID, member.RoomID, member.UserID)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MemberRepository) DeleteMember(member *entity.Member) error {
	query := dml.DeleteMemberQuery
	_, err := m.db.Exec(query, member.RoomID, member.UserID)
	if err != nil {
		return err
	}
	return nil
}
