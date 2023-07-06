package service

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"chat-server/utils"
)

type userService struct {
	repo use_case.UserStorage
}

func NewUserService(userRepo use_case.UserStorage) use_case.UserUseCase {
	return &userService{
		repo: userRepo,
	}
}

func (u *userService) CreateUser(req *entity.CreateUserReq) (*entity.CreateUserRes, error) {
	hashedPassword := utils.HashPassword(req.Password)
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := u.repo.CreateUser(user)
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
	return u.repo.SelectUserByEmailAndPassword(email, password)
}

func (u *userService) UserExists(id entity.ID) (bool, error) {
	_, err := u.repo.SelectUserByID(id)
	if err != nil {
		if err == use_case.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u *userService) EditUserProfile(req *entity.EditProfileReq) (*entity.EditProfileRes, error) {
	user, err := u.repo.SelectUserByID(req.ID)
	if err != nil {
		return nil, err
	}
	user.Username = req.Username

	updatedUser, err := u.repo.UpdateUser(user)
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
