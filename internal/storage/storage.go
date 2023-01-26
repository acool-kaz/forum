package storage

import (
	"context"
	"database/sql"
	"forum/internal/models"
)

const (
	sessionTable  = "sessions"
	userTable     = "users"
	postTable     = "posts"
	imageTable    = "images"
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
	GetOneBy(ctx context.Context) (models.User, error)
	Create(ctx context.Context, user models.User) (uint, error)
}

type Post interface {
	Create(ctx context.Context, post models.Post) (uint, error)
	GetAll(ctx context.Context) ([]models.FullPost, error)
	GetById(ctx context.Context, id uint) (models.FullPost, error)
	Delete(ctx context.Context, id uint) error
	SaveImages(ctx context.Context, postId uint, url string) error
}

type Tags interface {
	Create(ctx context.Context, tag models.Tag) error
}

type Comment interface {
	Create(ctx context.Context, comment models.Comment) error
	GetAll(ctx context.Context, postId uint) ([]models.FullComment, error)
}

type Reaction interface {
	CreateForPost(ctx context.Context, reaction models.Reaction) error
	CreateForComment(ctx context.Context, reaction models.Reaction) error
	Get(ctx context.Context, postId, commentId, userId uint) (models.Reaction, error)
	Change(ctx context.Context, id uint, newReact int) error
	Delete(ctx context.Context, id uint) error
}

type Storage struct {
	Session  Session
	User     User
	Post     Post
	Tags     Tags
	Comment  Comment
	Reaction Reaction
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Session:  newSessionStorage(db),
		User:     newUserStorage(db),
		Post:     newPostStorage(db),
		Tags:     newTagsStorage(db),
		Comment:  newCommentStorage(db),
		Reaction: newReactionStorage(db),
	}
}
