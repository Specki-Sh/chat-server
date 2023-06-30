package use_case

import (
	"chat-server/internal/domain/entity"
)

type UserUseCase interface {
	CreateUser(req *entity.CreateUserReq) (*entity.CreateUserRes, error)
	GetByEmailAndPassword(email entity.Email, password entity.HashPassword) (*entity.User, error)
	UserExists(id entity.ID) (bool, error)
	EditUserProfile(req *entity.EditProfileReq) (*entity.EditProfileRes, error)
}
