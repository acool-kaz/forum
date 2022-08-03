package service

import (
	"errors"
	"forum/internal/storage"
	"forum/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginNotFound   = errors.New("login not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUserName = errors.New("invalid username")
)

type Auth interface {
	CreateUser(user models.User) error
	GenerateSessionToken(login, password string) (string, time.Time, error)
	ParseSessionToken(token string) (models.User, error)
	DeleteSessionToken(token string) error
}

type AuthService struct {
	storage storage.Auth
}

func newAuthService(storage storage.Auth) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (s *AuthService) CreateUser(user models.User) error {
	var err error
	user.Username, err = checkUserName(user.Username)
	if err != nil {
		return err
	}
	user.Password, err = generateHashPassword(user.Password)
	if err != nil {
		return err
	}
	return s.storage.CreateUser(user)
}

func (s *AuthService) GenerateSessionToken(username, password string) (string, time.Time, error) {
	username, err := checkUserName(username)
	if err != nil {
		return "", time.Time{}, ErrLoginNotFound
	}
	user, err := s.storage.GetUserByLogin(username)
	if err != nil {
		return "", time.Time{}, ErrLoginNotFound
	}
	if err := compareHashAndPassword(user.Password, password); err != nil {
		return "", time.Time{}, ErrInvalidPassword
	}
	token := uuid.NewString()
	expiresAt := time.Now().Add(time.Hour * 12)
	if err := s.storage.SaveSessinToken(user.Username, token, expiresAt); err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (s *AuthService) ParseSessionToken(token string) (models.User, error) {
	user, err := s.storage.GetUserByToken(token)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *AuthService) DeleteSessionToken(token string) error {
	return s.storage.DeleteSessionToken(token)
}

func generateHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func compareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func checkUserName(username string) (string, error) {
	username = strings.ToLower(username)
	for _, char := range username {
		if char <= 32 || char >= 127 {
			return "", ErrInvalidUserName
		}
	}
	if len(username) >= 32 {
		return "", ErrInvalidUserName
	}
	return username, nil
}
