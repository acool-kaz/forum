package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/models"
)

type SessionStorage struct {
	db *sql.DB
}

func newSessionStorage(db *sql.DB) *SessionStorage {
	return &SessionStorage{
		db: db,
	}
}

func (s *SessionStorage) GetByToken(ctx context.Context, token string) (models.Session, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE token = $1;", sessionTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.Session{}, fmt.Errorf("session storage: get by token: %w", err)
	}
	defer prep.Close()

	var session models.Session
	if err = prep.QueryRowContext(ctx, token).Scan(&session.Id, &session.UserId, &session.Token, &session.ExpiresAt); err != nil {
		return models.Session{}, fmt.Errorf("session storage: get by token: %w", err)
	}

	return session, nil
}

func (s *SessionStorage) Create(ctx context.Context, session models.Session) error {
	query := fmt.Sprintf("INSERT INTO %s(user_id, token, expires_at) VALUES($1, $2, $3);", sessionTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("session storage: create: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, session.UserId, session.Token, session.ExpiresAt); err != nil {
		return fmt.Errorf("session storage: create: %w", err)
	}

	return nil
}

func (s *SessionStorage) Delete(ctx context.Context, token string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE token = $1", sessionTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("session storage: delete: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, token); err != nil {
		return fmt.Errorf("session storage: delete: %w", err)
	}

	return nil
}
