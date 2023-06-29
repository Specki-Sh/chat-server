package entity

import (
	"errors"
	"net/mail"
	"unicode"
)

// errors
var (
	ErrInvalidEmail   = errors.New("invalid email address")
	ErrPasswordLength = errors.New("password must be at least 8 characters long")
	ErrPasswordUpper  = errors.New("password must contain at least one uppercase letter")
	ErrPasswordLower  = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNumber = errors.New("password must contain at least one number")
)

type ValidateTypes interface {
	Validate() error
}

type Email string

func (e *Email) Validate() error {
	_, err := mail.ParseAddress(string(*e))
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

type Password string

// Validate checks if a password is valid.
// A password is considered valid if it meets the following criteria:
// - It is at least 8 characters long.
// - It contains at least one uppercase letter.
// - It contains at least one lowercase letter.
// - It contains at least one number.
func (p Password) Validate() error {
	if len(p) < 8 {
		return ErrPasswordLength
	}

	var hasUpper, hasLower, hasNumber bool
	for _, c := range p {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}

	if !hasUpper {
		return ErrPasswordUpper
	}
	if !hasLower {
		return ErrPasswordLower
	}
	if !hasNumber {
		return ErrPasswordNumber
	}

	return nil
}

type HashPassword string

type ID uint
