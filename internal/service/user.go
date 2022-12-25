package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStorage storage.User
}

func newUserService(userStorage storage.User) *UserService {
	return &UserService{
		userStorage: userStorage,
	}
}

func (s *UserService) GetById(ctx context.Context, id uint) (models.User, error) {
	user, err := s.userStorage.GetById(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("user service: get by id: %w", err)
		}
	}
	return user, nil
}

func (s *UserService) Create(ctx context.Context, user models.User) (uint, error) {
	var (
		err error
		id  uint
	)
	if err = validUser(user); err != nil {
		return 0, fmt.Errorf("user service: create: %w", err)
	}
	if _, err = s.userStorage.GetByUsername(ctx, user.Username); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user service: create: %w: username exist", models.ErrInvalidUser)
		}
	}
	if _, err = s.userStorage.GetByEmail(ctx, user.Email); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user service: create: %w: email exist", models.ErrInvalidUser)
		}
	}
	user.Password, err = generateHashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("user service: create: %w", err)
	}
	id, err = s.userStorage.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("user service: create: %w", err)
	}
	return id, nil
}

func generateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func validUser(user models.User) error {
	if len(user.Username) < 4 || len(user.Username) > 36 {
		return fmt.Errorf("%w: invalid username", models.ErrInvalidUser)
	}
	for _, char := range user.Username {
		if char <= 32 || char >= 127 {
			return fmt.Errorf("%w: invalid username", models.ErrInvalidUser)
		}
	}
	r := regexp.MustCompile(`[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	validEmail := r.MatchString(user.Email)
	if !validEmail {
		return fmt.Errorf("%w: invalid email", models.ErrInvalidUser)
	}
	if user.Password != user.VerifyPassword {
		return fmt.Errorf("%w: password dont match", models.ErrInvalidUser)
	}
	return nil
}
