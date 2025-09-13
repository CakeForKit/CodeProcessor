package models

import "github.com/google/uuid"

type Session struct {
	sessionID uuid.UUID
	userID    uuid.UUID
}

func NewSession(sessionID uuid.UUID, userID uuid.UUID) (Session, error) {
	return Session{
		sessionID: sessionID,
		userID:    userID,
	}, nil
}

func (s *Session) SessionID() uuid.UUID {
	return s.sessionID
}

func (s *Session) UserID() uuid.UUID {
	return s.userID
}
