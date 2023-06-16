package use_case

import (
	"chat-server/internal/domain/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByEmailAndPassword(email string, password string) (*entity.User, error)
}
