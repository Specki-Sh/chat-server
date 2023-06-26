package repository

import (
	"chat-server/internal/domain/entity"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	var id int
	query := "INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := u.db.QueryRow(query, user.Username, user.Password, user.Email).Scan(&id)
	if err != nil {
		return &entity.User{}, err
	}

	user.ID = id
	return user, nil
}

func (u *UserRepository) GetUserByEmailAndPassword(email string, password string) (*entity.User, error) {
	var user entity.User
	query := "SELECT id, username, password, email FROM users WHERE email = $1 AND password = $2"
	err := u.db.QueryRow(query, email, password).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

func (u *UserRepository) SelectUserByID(id int) (*entity.User, error) {
	var user entity.User
	query := "SELECT id, username, password, email FROM users WHERE id = $1"
	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}
