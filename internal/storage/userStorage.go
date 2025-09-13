package storage

import (
	"codeProcessor/internal/models"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UserStorage interface {
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetByLogin(login string) (*models.User, error)
	AddUser(user *models.User) error
}

func NewUserStorage() (UserStorage, error) {
	return &userStorage{
		users: make(map[uuid.UUID]*models.User),
	}, nil
}

type userStorage struct {
	users map[uuid.UUID]*models.User
}

var (
	ErrNoUser            = errors.New("no user in map")
	ErrUserIDAlreadExist = errors.New("user id already exist")
)

func (ts *userStorage) GetUserByID(id uuid.UUID) (*models.User, error) {
	if t, ok := ts.users[id]; ok {
		return t, nil
	} else {
		return nil, ErrNoUser
	}
}

func (ts *userStorage) GetByLogin(login string) (*models.User, error) {
	fmt.Print(ts.users)
	for _, v := range ts.users {
		if v.Login() == login {
			return v, nil
		}
	}
	return nil, ErrNoUser
}

func (ts *userStorage) AddUser(user *models.User) error {
	if _, ok := ts.users[user.ID()]; ok {
		return ErrUserIDAlreadExist
	}
	ts.users[user.ID()] = user
	return nil
}
