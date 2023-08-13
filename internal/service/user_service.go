package service

import (
	"context"
	"fmt"

	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"chat-server/utils"
)

type userService struct {
	userRepo  use_case.UserStorage
	userCache use_case.UserCacheStorage
}

func NewUserService(
	userRepo use_case.UserStorage,
	userCache use_case.UserCacheStorage,
) use_case.UserUseCase {
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
		return nil, fmt.Errorf("userService.CreateUser: %w", err)
	}

	res := &entity.CreateUserRes{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (u *userService) GetByEmailAndPassword(
	email entity.Email,
	password entity.HashPassword,
) (*entity.User, error) {
	user, err := u.userRepo.SelectUserByEmailAndPassword(email, password)
	if err != nil {
		return nil, fmt.Errorf("userService.GetByEmailAndPassword: %w", err)
	}
	return user, nil
}

func (u *userService) UserExists(id entity.ID) (bool, error) {
	_, err := u.userRepo.SelectUserByID(id)
	if err != nil {
		if err == use_case.ErrUserNotFound {
			return false, nil
		}
		return false, fmt.Errorf("userService.UserExists: %w", err)
	}
	return true, nil
}

func (u *userService) EditUserProfile(req *entity.EditProfileReq) (*entity.EditProfileRes, error) {
	user, err := u.userRepo.SelectUserByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("userService.EditUserProfile: %w", err)
	}
	user.Username = req.Username

	updatedUser, err := u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("userService.EditUserProfile: %w", err)
	}

	res := &entity.EditProfileRes{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
	}

	return res, nil
}

func (u *userService) StoreUserData(
	ctx context.Context,
	secretCode string,
	userData *entity.UserData,
) error {
	if err := u.userCache.SetUserData(ctx, secretCode, userData); err != nil {
		return fmt.Errorf("userService.StoreUserData: %w", err)
	}
	return nil
}

func (u *userService) RetrieveUserData(
	ctx context.Context,
	secretCode string,
) (*entity.UserData, error) {
	userData, err := u.userCache.GetUserData(ctx, secretCode)
	if err != nil {
		return nil, fmt.Errorf("userService.RetrieveUserData: %w", err)
	}
	if err := u.userCache.DeleteUserData(ctx, secretCode); err != nil {
		return nil, fmt.Errorf("userService.RetrieveUserData: %w", err)
	}
	return userData, nil
}
