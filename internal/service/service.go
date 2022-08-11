package service

import (
	"forum/internal/storage"
)

type Service struct {
	Auth
	Post
	Comment
	LikeDislikePost
	LikeDislikeComment
	User
	Notify
}

func NewService(storages *storage.Storage) *Service {
	return &Service{
		Auth:               newAuthService(storages.Auth),
		Post:               newPostService(storages.Post),
		Comment:            newCommentService(storages.Comment),
		LikeDislikePost:    newLikeDislikeService(storages.LikeDislikePost),
		LikeDislikeComment: newLikeDislikeCommentService(storages.LikeDislikeComment),
		User:               newUserService(storages.User),
		Notify:             newNotifyService(storages.Notify),
	}
}
