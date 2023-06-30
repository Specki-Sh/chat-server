package repository

import (
	"chat-server/internal/domain/entity"
	"chat-server/internal/domain/use_case"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return nil, use_case.ErrUserInvalid
	}
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := u.db.QueryRow(query, user.Username, user.Password, user.Email).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) SelectUserByEmailAndPassword(email entity.Email, password entity.HashPassword) (*entity.User, error) {
	var user entity.User
	query := "SELECT id, username, password, email FROM users WHERE email = $1 AND password = $2"
	err := u.db.QueryRow(query, email, password).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, use_case.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) SelectUserByID(id entity.ID) (*entity.User, error) {
	var user entity.User
	query := "SELECT id, username, password, email FROM users WHERE id = $1"
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, use_case.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	query := `
        UPDATE users
        SET username = $1, password = $2, email = $3
        WHERE id = $4
    `
	_, err := u.db.Exec(query, user.Username, user.Password, user.Email, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
