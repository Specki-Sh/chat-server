package entity

type User struct {
	ID       int          `json:"id"`
	Username string       `json:"username"`
	Email    Email        `json:"email"`
	Password HashPassword `json:"password"`
}

type CreateUserReq struct {
	Username string   `json:"username"`
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

type CreateUserRes struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    Email  `json:"email"`
}

type EditProfileReq struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type EditProfileRes struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    Email  `json:"email"`
}
