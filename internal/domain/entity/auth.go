package entity

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRes struct {
	AccessToken string `json:"-"`
	ID          int    `json:"id"`
	Username    string `json:"username"`
}
