package storage

import (
	"context"
	"database/sql"
	"forum/models"
)

const (
	sessionTable  = "sessions"
	userTable     = "users"
	postTable     = "posts"
	tagTable      = "tags"
	commentTable  = "comments"
	reactionTable = "reactions"
)

type Session interface {
	GetByToken(ctx context.Context, token string) (models.Session, error)
	Create(ctx context.Context, session models.Session) error
	Delete(ctx context.Context, token string) error
}

type User interface {
	GetByUsername(ctx context.Context, username string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetById(ctx context.Context, id uint) (models.User, error)
	Create(ctx context.Context, user models.User) (uint, error)
}

type Post interface {
	Create(ctx context.Context, post models.Post) (uint, error)
	GetAll(ctx context.Context) ([]models.FullPost, error)
	GetById(ctx context.Context, id uint) (models.FullPost, error)
	Delete(ctx context.Context, id uint) error
}

type Tags interface {
	Create(ctx context.Context, tag models.Tag) error
}

type Comment interface {
	Create(ctx context.Context, comment models.Comment) error
	GetAll(ctx context.Context, postId uint) ([]models.FullComment, error)
}

type Storage struct {
	Session
	User
	Post
	Tags
	Comment
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Session: newSessionStorage(db),
		User:    newUserStorage(db),
		Post:    newPostStorage(db),
		Tags:    newTagsStorage(db),
		Comment: newCommentStorage(db),
	}
}
