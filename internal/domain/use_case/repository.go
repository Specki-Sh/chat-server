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
	SelectUserByEmailAndPassword(email entity.Email, password entity.HashPassword) (*entity.User, error)
	SelectUserByID(id entity.ID) (*entity.User, error)
}

type RoomRepository interface {
	InsertRoom(room *entity.Room) (*entity.Room, error)
	SelectRoomByID(id entity.ID) (*entity.Room, error)
	UpdateRoom(room *entity.Room) error
	DeleteRoom(id entity.ID) error
}

type MessageRepository interface {
	InsertMessage(message *entity.Message) (*entity.Message, error)
	SelectMessage(id entity.ID) (*entity.Message, error)
	UpdateMessage(message *entity.Message) error
	SoftDeleteMessageByID(id entity.ID) error
	SoftDeleteMessagesByRoomID(roomID entity.ID) error

	SelectMessagePaginate(roomID entity.ID, perPage uint, page uint) ([]*entity.Message, error)
	SelectMessagesPaginateReverse(roomID entity.ID, perPage uint, page uint) ([]*entity.Message, error)
}

type MemberRepository interface {
	InsertMember(member *entity.Member) (*entity.Member, error)
	SelectMembersByRoomID(roomID entity.ID) ([]*entity.Member, error)
	UpdateMember(member *entity.Member) (*entity.Member, error)
	DeleteMember(member *entity.Member) error
}
