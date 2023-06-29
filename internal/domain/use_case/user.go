package use_case

import (
	"chat-server/internal/domain/entity"
)

type UserUseCase interface {
	CreateUser(req *entity.CreateUserReq) (*entity.CreateUserRes, error)
	GetByEmailAndPassword(email string, password string) (*entity.User, error)
	UserExists(id int) (bool, error)
	EditUserProfile(req *entity.EditProfileReq) (*entity.EditProfileRes, error)
}
