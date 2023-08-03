package service

import (
	"context"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"chat-server/utils"
)

type userService struct {
	userRepo  use_case.UserStorage
	userCache use_case.UserCacheStorage
}

func NewUserService(userRepo use_case.UserStorage, userCache use_case.UserCacheStorage) use_case.UserUseCase {
	return &userService{
		userRepo:  userRepo,
		userCache: userCache,
	}
}

func (u *userService) CreateUser(req *entity.CreateUserReq) (*entity.CreateUserRes, error) {
	hashedPassword := utils.HashPassword(req.Password)
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	res := &entity.CreateUserRes{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (u *userService) GetByEmailAndPassword(email entity.Email, password entity.HashPassword) (*entity.User, error) {
	return u.userRepo.SelectUserByEmailAndPassword(email, password)
}

func (u *userService) UserExists(id entity.ID) (bool, error) {
	_, err := u.userRepo.SelectUserByID(id)
	if err != nil {
		if err == use_case.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u *userService) EditUserProfile(req *entity.EditProfileReq) (*entity.EditProfileRes, error) {
	user, err := u.userRepo.SelectUserByID(req.ID)
	if err != nil {
		return nil, err
	}
	user.Username = req.Username

	updatedUser, err := u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	res := &entity.EditProfileRes{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
	}

	return res, nil
}

func (u *userService) StoreUserData(ctx context.Context, secretCode string, userData *entity.UserData) error {
	return u.userCache.SetUserData(ctx, secretCode, userData)
}

func (u *userService) RetrieveUserData(ctx context.Context, secretCode string) (*entity.UserData, error) {
	userData, err := u.userCache.GetUserData(ctx, secretCode)
	if err != nil {
		return nil, err
	}
	if err := u.userCache.DeleteUserData(ctx, secretCode); err != nil {
		return nil, err
	}
	return userData, nil
}
