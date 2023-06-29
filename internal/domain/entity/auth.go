package entity

type SignInReq struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

type SignInRes struct {
	AccessToken string `json:"-"`
	ID          ID     `json:"id"`
	Username    string `json:"username"`
}
