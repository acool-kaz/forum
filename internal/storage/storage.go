package storage

import (
	"database/sql"
)

type Storage struct {
	Auth
	Post
	Comment
	LikeDislikePost
	LikeDislikeComment
	User
	Notify
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Auth:               newAuthStorage(db),
		Post:               newPostStorage(db),
		Comment:            newCommentStorage(db),
		LikeDislikePost:    newLikeDislikePostStorage(db),
		LikeDislikeComment: newLikeDislikeCommentStorage(db),
		User:               newUserStorage(db),
		Notify:             newNotifyStorage(db),
	}
}
