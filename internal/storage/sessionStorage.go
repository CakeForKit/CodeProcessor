package storage

import (
	"codeProcessor/internal/models"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type SessionStorage interface {
	Init(sid uuid.UUID, uid uuid.UUID) (models.Session, error) // implements the initialization of a session, and returns a new session if it succeeds.
	Read(sid uuid.UUID) (models.Session, error)                // returns a session represented by the corresponding sid. Return error if it does not already exist.
	Destroy(sid uuid.UUID) error                               // given an sid, deletes the corresponding session.
	// GC(maxLifeTime int64)                                      // (garbage collection) deletes expired session variables according to maxLifeTime. ()
}

var (
	ErrSessionStorage       = errors.New("SessionStorage ")
	ErrSessionAlreadyExists = errors.New("SessionStorage: session with this ID already exists")
	ErrSessionNotFound      = errors.New("SessionStorage: session not found")
	// ErrStorageUnavailable   = errors.New("SessionStorage: failed to read from storage")s
)

var SessionStorages = make(map[string]SessionStorage)

func RegisterSessionStorage(name string, storage SessionStorage) {
	if storage == nil {
		panic("session: Register SessionStorage is nil")
	}
	if _, dup := SessionStorages[name]; dup {
		panic("session: Register called twice for SessionStorage " + name)
	}
	SessionStorages[name] = storage
}

func RegisterAllSessionStorages() {
	SessionStorages["mem"] = &memSessionStorage{st: make(map[uuid.UUID]models.Session)}
}

type memSessionStorage struct {
	lock sync.Mutex
	st   map[uuid.UUID]models.Session
}

func (s *memSessionStorage) Init(sid uuid.UUID, uid uuid.UUID) (models.Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.st[sid]; ok {
		return models.Session{}, fmt.Errorf("%w", ErrSessionAlreadyExists)
	}

	sess, err := models.NewSession(sid, uid)
	if err != nil {
		return models.Session{}, fmt.Errorf("%w - %w", ErrSessionStorage, err)
	}
	s.st[sid] = sess
	return sess, nil
}

func (s *memSessionStorage) Read(sid uuid.UUID) (models.Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if sess, ok := s.st[sid]; ok {
		return sess, nil
	}
	return models.Session{}, fmt.Errorf("%w", ErrSessionNotFound)
}

func (s *memSessionStorage) Destroy(sid uuid.UUID) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.st[sid]; !ok {
		return fmt.Errorf("%w", ErrSessionNotFound)
	}
	delete(s.st, sid)
	return nil
}
