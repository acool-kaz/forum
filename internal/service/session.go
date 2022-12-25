package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SessionService struct {
	sessionStorage storage.Session
	userStorage    storage.User
}

func newSessionService(sessionStorage storage.Session, userStorage storage.User) *SessionService {
	return &SessionService{
		sessionStorage: sessionStorage,
		userStorage:    userStorage,
	}
}

func (s *SessionService) GenerateSessionToken(ctx context.Context, username, password string) (models.Session, error) {
	var (
		err      error
		tempUser models.User
		session  models.Session
	)
	tempUser, err = s.userStorage.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, fmt.Errorf("session service: generate session token: %w", models.ErrUserNotFound)
		}
		return models.Session{}, fmt.Errorf("session service: generate session token: %w", err)
	}
	if err = compareHashAndPassword(tempUser.Password, password); err != nil {
		return models.Session{}, fmt.Errorf("session service: generate session token: %w", models.ErrUserNotFound)
	}
	session = models.Session{
		UserId:    tempUser.Id,
		Token:     uuid.NewString(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
	}
	if err = s.sessionStorage.Create(ctx, session); err != nil {
		return models.Session{}, fmt.Errorf("session service: generate session token: %w", err)
	}
	return session, nil
}

func (s *SessionService) ParseSessionToken(ctx context.Context, token string) (models.Session, error) {
	session, err := s.sessionStorage.GetByToken(ctx, token)
	if err != nil {
		return models.Session{}, fmt.Errorf("session service: parse session token: %w", err)
	}
	return session, nil
}

func (s *SessionService) DeleteSessionToken(ctx context.Context, token string) error {
	if err := s.sessionStorage.Delete(ctx, token); err != nil {
		return fmt.Errorf("session service: delete session token: %w", err)
	}
	return nil
}

func compareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
