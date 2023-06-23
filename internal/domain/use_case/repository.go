package use_case

import (
	"chat-server/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByEmailAndPassword(email string, password string) (*entity.User, error)
}

type RoomRepository interface {
	InsertRoom(room *entity.Room) (*entity.Room, error)
	SelectRoom(id int) (*entity.Room, error)
	UpdateRoom(room *entity.Room) error
	DeleteRoom(id int) error
}

type MessageRepository interface {
	InsertMessage(message *entity.Message) (*entity.Message, error)
	SelectMessage(id int) (*entity.Message, error)
	UpdateMessage(message *entity.Message) error
	DeleteMessage(id int) error

	SelectMessagePaginate(perPage int, page int) ([]*entity.Message, error)
}
