package service

import (
	"forum/internal/storage"
	"forum/models"
)

type Comment interface {
	GetComments(postId int) ([]models.Comment, error)
	CreateComment(comment models.Comment) error
	GetPostIdByCommentId(commentId int) (int, error)
}

type CommentService struct {
	storage storage.Comment
}

func newCommentService(storage storage.Comment) *CommentService {
	return &CommentService{
		storage: storage,
	}
}

func (s *CommentService) GetComments(postId int) ([]models.Comment, error) {
	return s.storage.GetComments(postId)
}

func (s *CommentService) CreateComment(comment models.Comment) error {
	return s.storage.CreateComment(comment)
}

func (s *CommentService) GetPostIdByCommentId(commentId int) (int, error) {
	return s.storage.GetPostIdByCommentId(commentId)
}
