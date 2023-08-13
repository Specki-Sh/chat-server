package use_case

import (
	"context"

	"chat-server/internal/domain/entity"
)

type UserUseCase interface {
	CreateUser(req *entity.CreateUserReq) (*entity.CreateUserRes, error)
	GetByEmailAndPassword(email entity.Email, password entity.HashPassword) (*entity.User, error)
	UserExists(id entity.ID) (bool, error)
	EditUserProfile(req *entity.EditProfileReq) (*entity.EditProfileRes, error)

	StoreUserData(ctx context.Context, secretCode string, userData *entity.UserData) error
	RetrieveUserData(ctx context.Context, secretCode string) (*entity.UserData, error)
}
