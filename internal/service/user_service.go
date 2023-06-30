package service

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"chat-server/internal/repository"
	"chat-server/utils"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(storage *repository.UserRepository) *UserService {
	return &UserService{storage}
}

func (u *UserService) CreateUser(req *entity.CreateUserReq) (*entity.CreateUserRes, error) {
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

func (u *UserService) GetByEmailAndPassword(email entity.Email, password entity.HashPassword) (*entity.User, error) {
	return u.repo.SelectUserByEmailAndPassword(email, password)
}

func (u *UserService) UserExists(id entity.ID) (bool, error) {
	_, err := u.repo.SelectUserByID(id)
	if err != nil {
		if err == use_case.ErrUserNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (u *UserService) EditUserProfile(req *entity.EditProfileReq) (*entity.EditProfileRes, error) {
	return nil, nil
}
