package entity

type SignInReq struct {
	Email    Email    `json:"email"`
	Password Password `json:"password"`
}

func (s *SignInReq) Validate() error {
	if err := s.Email.Validate(); err != nil {
		return err
	}
	if err := s.Password.Validate(); err != nil {
		return err
	}
	return nil
}

type SignInRes struct {
	AccessToken string `json:"-"`
	ID          ID     `json:"id"`
	Username    string `json:"username"`
}
