package use_case

import (
	"chat-server/internal/domain/entity"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserInvalid  = errors.New("user data is invalid or incomplete")
	ErrRoomNotFound = errors.New("room not found")
	ErrRoomInvalid  = errors.New("room data is invalid or incomplete")
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	SelectUserByEmailAndPassword(email string, password string) (*entity.User, error)
	SelectUserByID(id int) (*entity.User, error)
}

type RoomRepository interface {
	InsertRoom(room *entity.Room) (*entity.Room, error)
	SelectRoomByID(id int) (*entity.Room, error)
	UpdateRoom(room *entity.Room) error
	DeleteRoom(id int) error
}

type MessageRepository interface {
	InsertMessage(message *entity.Message) (*entity.Message, error)
	SelectMessage(id int) (*entity.Message, error)
	UpdateMessage(message *entity.Message) error
	SoftDeleteMessageByID(id int) error
	SoftDeleteMessagesByRoomID(roomID int) error

	SelectMessagePaginate(roomID int, perPage int, page int) ([]*entity.Message, error)
	SelectMessagesPaginateReverse(roomID int, perPage int, page int) ([]*entity.Message, error)
}

type MemberRepository interface {
	InsertMember(member *entity.Member) (*entity.Member, error)
	SelectMembersByRoomID(roomID int) ([]*entity.Member, error)
	UpdateMember(member *entity.Member) (*entity.Member, error)
	DeleteMember(member *entity.Member) error
}
