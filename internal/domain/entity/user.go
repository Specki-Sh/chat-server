package entity

type User struct {
	ID       ID             `json:"id"`
	Username NonEmptyString `json:"username"`
	Email    Email          `json:"email"`
	Password HashPassword   `json:"password"`
}

type CreateUserReq struct {
	Username NonEmptyString `json:"username"`
	Email    Email          `json:"email"`
	Password Password       `json:"password"`
}

func (c *CreateUserReq) Validate() error {
	if err := c.Username.Validate(); err != nil {
		return err
	}
	if err := c.Email.Validate(); err != nil {
		return err
	}
	if err := c.Password.Validate(); err != nil {
		return err
	}
	return nil
}

type CreateUserRes struct {
	ID       ID             `json:"id"`
	Username NonEmptyString `json:"username"`
	Email    Email          `json:"email"`
}

type EditProfileReq struct {
	ID       ID             `json:"id"`
	Username NonEmptyString `json:"username"`
}

func (e *EditProfileReq) Validate() error {
	if err := e.ID.Validate(); err != nil {
		return err
	}
	if err := e.Username.Validate(); err != nil {
		return err
	}
	return nil
}

type EditProfileRes struct {
	ID       ID             `json:"id"`
	Username NonEmptyString `json:"username"`
	Email    Email          `json:"email"`
}

type UserData struct {
	Username NonEmptyString `json:"username"`
	Email    Email          `json:"email"`
	Password HashPassword   `json:"password"`
}
