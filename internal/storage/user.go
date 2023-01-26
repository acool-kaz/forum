package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/internal/models"
	"strings"
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

func (s *UserStorage) GetOneBy(ctx context.Context) (models.User, error) {
	args := []interface{}{}
	argsStr := []string{}
	argsNum := 1

	userId := ctx.Value(models.UserId)
	if userId != nil {
		argsStr = append(argsStr, fmt.Sprintf("id = $%d", argsNum))
		args = append(args, userId.(uint))
		argsNum++
	}

	username := ctx.Value(models.Username)
	if username != nil {
		argsStr = append(argsStr, fmt.Sprintf("username = $%d", argsNum))
		args = append(args, username.(string))
		argsNum++
	}

	email := ctx.Value(models.Email)
	if email != nil {
		argsStr = append(argsStr, fmt.Sprintf("email = $%d", argsNum))
		args = append(args, email.(string))
		argsNum++
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s;", userTable, strings.Join(argsStr, " AND "))

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.User{}, fmt.Errorf("user storage: get on by: %w", err)
	}
	defer prep.Close()

	var user models.User
	if err = prep.QueryRowContext(ctx, args...).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return models.User{}, fmt.Errorf("user storage: get on by: %w", err)
	}

	return user, nil
}
