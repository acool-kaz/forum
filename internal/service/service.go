package service

import (
	"context"
	"forum/internal/storage"
	"forum/models"
)

type Session interface {
	GenerateSessionToken(ctx context.Context, username, password string) (models.Session, error)
	ParseSessionToken(ctx context.Context, token string) (models.Session, error)
	DeleteSessionToken(ctx context.Context, token string) error
}

type User interface {
	GetById(ctx context.Context, id uint) (models.User, error)
	Create(ctx context.Context, user models.User) (uint, error)
}

type Post interface {
	Create(ctx context.Context, post models.Post) (uint, error)
	GetAll(ctx context.Context) ([]models.FullPost, error)
	GetById(ctx context.Context, id uint) (models.FullPost, error)
}

type Comment interface {
	Create(ctx context.Context, comment models.Comment) error
}

type Service struct {
	Session
	User
	Post
	Comment
}

func NewService(storages *storage.Storage) *Service {
	return &Service{
		Session: newSessionService(storages.Session, storages.User),
		User:    newUserService(storages.User),
		Post:    newPostService(storages.Post, storages.Tags, storages.Comment),
		Comment: newCommentService(storages.Comment),
	}
}
