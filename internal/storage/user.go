package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/models"
)

type UserStorage struct {
	db *sql.DB
}

func newUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Create(ctx context.Context, user models.User) (uint, error) {
	query := fmt.Sprintf("INSERT INTO %s(username, email, password) VALUES($1, $2, $3) RETURNING id;", userTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("user storage: create: %w", err)
	}
	defer prep.Close()

	var id uint
	if err = prep.QueryRowContext(ctx, user.Username, user.Email, user.Password).Scan(&id); err != nil {
		return 0, fmt.Errorf("user storage: create: %w", err)
	}

	return id, nil
}

func (s *UserStorage) GetById(ctx context.Context, id uint) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.User{}, fmt.Errorf("user storage: get by id: %w", err)
	}
	defer prep.Close()

	var user models.User
	if err = prep.QueryRowContext(ctx, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return models.User{}, fmt.Errorf("user storage: get by id: %w", err)
	}

	return user, nil
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", userTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.User{}, fmt.Errorf("user storage: get by username: %w", err)
	}
	defer prep.Close()

	var user models.User
	if err = prep.QueryRowContext(ctx, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return models.User{}, fmt.Errorf("user storage: get by username: %w", err)
	}

	return user, nil
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1", userTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.User{}, fmt.Errorf("user storage: get by email: %w", err)
	}
	defer prep.Close()

	var user models.User
	if err = prep.QueryRowContext(ctx, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return models.User{}, fmt.Errorf("user storage: get by email: %w", err)
	}

	return user, nil
}
