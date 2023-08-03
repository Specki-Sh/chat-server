package use_case

import (
	"context"
	"errors"
	"time"

	"chat-server/internal/domain/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserInvalid  = errors.New("user data is invalid or incomplete")
	ErrRoomNotFound = errors.New("room not found")
	ErrRoomInvalid  = errors.New("room data is invalid or incomplete")
)

type UserStorage interface {
	CreateUser(user *entity.User) (*entity.User, error)
	SelectUserByEmailAndPassword(email entity.Email, password entity.HashPassword) (*entity.User, error)
	SelectUserByID(id entity.ID) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
}

type UserCacheStorage interface {
	SetUserData(ctx context.Context, secretCode string, userData *entity.UserData) error
	GetUserData(ctx context.Context, secretCode string) (*entity.UserData, error)
	DeleteUserData(ctx context.Context, secretCode string) error
}

type RoomStorage interface {
	InsertRoom(room *entity.Room) (*entity.Room, error)
	SelectRoomByID(id entity.ID) (*entity.Room, error)
	UpdateRoom(room *entity.Room) error
	DeleteRoom(id entity.ID) error
}

type MessageStorage interface {
	InsertMessage(message *entity.Message) (*entity.Message, error)
	SelectMessage(id entity.ID) (*entity.Message, error)
	UpdateMessage(message *entity.Message) error
	SoftDeleteMessageByID(id entity.ID) error
	SoftDeleteMessageBulkByRoomID(roomID entity.ID) error

	SelectMessageBulkPaginate(roomID entity.ID, perPage uint, page uint) ([]entity.Message, error)
	SelectMessageBulkPaginateReverse(roomID entity.ID, perPage uint, page uint) ([]entity.Message, error)
}

type MemberStorage interface {
	InsertMember(member *entity.Member) (*entity.Member, error)
	SelectMemberBulkByRoomID(roomID entity.ID) ([]entity.Member, error)
	UpdateMember(member *entity.Member) (*entity.Member, error)
	DeleteMember(member *entity.Member) error
}

type TokenStorage interface {
	SetInvalidRefreshToken(ctx context.Context, userID entity.ID, refreshToken string, expiresIn time.Duration) error
	InvalidRefreshTokenExists(ctx context.Context, refreshToken string) (bool, error)
}
