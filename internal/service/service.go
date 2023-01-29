package service

import (
	"context"
	"forum/internal/config"
	"forum/internal/models"
	"forum/internal/storage"
	"mime/multipart"
)

type Session interface {
	GenerateSessionToken(ctx context.Context, username, password string) (models.Session, error)
	ParseSessionToken(ctx context.Context, token string) (models.Session, error)
	DeleteSessionToken(ctx context.Context, token string) error
}

type User interface {
	GetOneBy(ctx context.Context) (models.User, error)
	Create(ctx context.Context, user models.User) (uint, error)
}

type Post interface {
	Create(ctx context.Context, post models.Post, files []*multipart.FileHeader) (uint, error)
	GetAll(ctx context.Context) ([]models.FullPost, error)
	GetById(ctx context.Context, id uint) (models.FullPost, error)
}

type Comment interface {
	Create(ctx context.Context, comment models.Comment) error
}

type Reaction interface {
	Set(ctx context.Context, postId, commentId string, react int, userId uint) error
}

type Service struct {
	Session  Session
	User     User
	Post     Post
	Comment  Comment
	Reaction Reaction
}

func NewService(storages *storage.Storage, cfg *config.Config) *Service {
	return &Service{
		Session:  newSessionService(storages.Session, storages.User),
		User:     newUserService(storages.User),
		Post:     newPostService(storages.Post, storages.Tags, storages.Comment, storages.Image, cfg),
		Comment:  newCommentService(storages.Comment),
		Reaction: newReactionService(storages.Reaction),
	}
}
