package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	id              uuid.UUID
	login           string
	hashed_password string
}

var (
	ErrEmptyLogin         = errors.New("login cannot be empty")
	ErrEmptyPassword      = errors.New("password cannot be empty")
	ErrInvalidLoginLength = errors.New("login must be between 3 and 50 characters")
	ErrValidateUser       = errors.New("not valid datas user's")
)

func NewUser(id uuid.UUID, login string, hashed_password string) (User, error) {
	user := User{
		id:              id,
		login:           login,
		hashed_password: hashed_password,
	}
	if !user.Validate() {
		return User{}, ErrValidateUser
	}
	return user, nil
}

func (u *User) Validate() bool {
	if strings.TrimSpace(u.login) == "" ||
		strings.TrimSpace(u.hashed_password) == "" ||
		len(u.login) < 3 || len(u.login) > 50 {
		return false
	}
	return true
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Login() string {
	return u.login
}

func (u *User) HashedPassword() string {
	return u.hashed_password
}

func (u *User) SetID(id uuid.UUID) {
	u.id = id
}

func (u *User) SetLogin(login string) error {
	if strings.TrimSpace(login) == "" {
		return ErrEmptyLogin
	}
	if len(login) < 3 || len(login) > 50 {
		return ErrInvalidLoginLength
	}
	u.login = login
	return nil
}

func (u *User) SetHashedPassword(hashed_password string) error {
	if strings.TrimSpace(hashed_password) == "" {
		return ErrEmptyPassword
	}
	u.hashed_password = hashed_password
	return nil
}

func (u *User) SetLoginWithValidation(login string) error {
	return u.SetLogin(login)
}

func (u *User) SetHashedPasswordWithValidation(hashed_password string) error {
	return u.SetHashedPassword(hashed_password)
}
