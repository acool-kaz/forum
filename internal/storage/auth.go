package storage

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"
)

type Auth interface {
	CreateUser(user models.User) error
	GetUserByLogin(login string) (models.User, error)
	SaveSessinToken(login, token string, expiresAt time.Time) error
	GetUserByToken(token string) (models.User, error)
	DeleteSessionToken(token string) error
}

type AuthStorage struct {
	db *sql.DB
}

func newAuthStorage(db *sql.DB) *AuthStorage {
	return &AuthStorage{
		db: db,
	}
}

func (s *AuthStorage) CreateUser(user models.User) error {
	query := `INSERT INTO user(email, username, hashPassword) VALUES ($1, $2, $3);`
	_, err := s.db.Exec(query, user.Email, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("storage: create user: %w", err)
	}
	return nil
}

func (s *AuthStorage) GetUserByLogin(username string) (models.User, error) {
	query := `SELECT id, email, username, hashPassword FROM user WHERE username=$1;`
	row := s.db.QueryRow(query, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by login: %w", err)
	}
	return user, nil
}

func (s *AuthStorage) SaveSessinToken(username, token string, expiresAt time.Time) error {
	query := `UPDATE user SET session_token = $1, expiresAt = $2 WHERE username = $3;`
	_, err := s.db.Exec(query, token, expiresAt, username)
	if err != nil {
		return fmt.Errorf("storage: save session token: %w", err)
	}
	return nil
}

func (s *AuthStorage) GetUserByToken(token string) (models.User, error) {
	query := `SELECT id, email, username, hashPassword, expiresAt FROM user WHERE session_token=$1;`
	row := s.db.QueryRow(query, token)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.ExpiresAt)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by token: %w", err)
	}
	return user, nil
}

func (s *AuthStorage) DeleteSessionToken(token string) error {
	query := `UPDATE user SET session_token = NULL, expiresAt = NULL WHERE session_token = $1;`
	_, err := s.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("storage: delete session token: %w", err)
	}
	return nil
}
