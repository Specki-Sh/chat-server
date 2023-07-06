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
	TokenPair TokenPair      `json:"token_pair"`
	ID        ID             `json:"id"`
	Username  NonEmptyString `json:"username"`
}

type RefreshTokenReq struct {
	RefreshToken string         `json:"refresh_token"`
	ID           ID             `json:"id"`
	Username     NonEmptyString `json:"username"`
}

type RefreshTokenRes struct {
	TokenPair TokenPair      `json:"token_pair"`
	ID        ID             `json:"id"`
	Username  NonEmptyString `json:"username"`
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token"`
}
