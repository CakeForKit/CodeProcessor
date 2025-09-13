package auth

import (
	"codeProcessor/internal/models"
	jsonrep "codeProcessor/internal/models/jsonRep"
	"codeProcessor/internal/services/hasher"
	"codeProcessor/internal/storage"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type AuthUserServ interface {
	RegisterUser(authUser *jsonrep.UserAuth) error
	LoginUser(authUser *jsonrep.UserAuth) (token uuid.UUID, err error)
	CheckSession(sid uuid.UUID) bool
}

var (
	ErrAuthServ = errors.New("AuthServ ")
	ErrPassword = errors.New("AuthServ: wrong password")
)

type authUserServ struct {
	sessionStorage storage.SessionStorage
	hashServ       hasher.Hasher
	userStorage    storage.UserStorage
	// maxlifetime int64
}

func NewAuthUserServ(storage storage.SessionStorage, hashServ hasher.Hasher, userStorage storage.UserStorage) (AuthUserServ, error) {
	return &authUserServ{
		sessionStorage: storage,
		hashServ:       hashServ,
		userStorage:    userStorage,
	}, nil
}

func (s *authUserServ) RegisterUser(authUser *jsonrep.UserAuth) error {
	hashed_password, err := s.hashServ.HashPassword(authUser.Password)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrAuthServ, err)
	}
	user, err := models.NewUser(uuid.New(), authUser.Login, hashed_password)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrAuthServ, err)
	}
	err = s.userStorage.AddUser(&user)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrAuthServ, err)
	}
	_, err = s.sessionStorage.Init(uuid.New(), user.ID())
	if err != nil {
		return fmt.Errorf("%w - %w", ErrAuthServ, err)
	}
	return nil
}

func (s *authUserServ) LoginUser(authUser *jsonrep.UserAuth) (token uuid.UUID, err error) {
	user, err := s.userStorage.GetByLogin(authUser.Login)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w - %w", ErrAuthServ, err)
	}

	err = s.hashServ.CheckPassword(authUser.Password, user.HashedPassword())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w - %w", ErrAuthServ, ErrPassword)
	}

	session, err := s.sessionStorage.Init(uuid.New(), user.ID())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w - %w", ErrAuthServ, err)
	}
	token = session.SessionID()
	return
}

func (s *authUserServ) CheckSession(sid uuid.UUID) bool {
	_, err := s.sessionStorage.Read(sid)
	return err == nil
}
