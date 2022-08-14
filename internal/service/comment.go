package service

import (
	"errors"
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"strings"
)

var ErrInvalidComment = errors.New("invalid comment")

type Comment interface {
	GetComments(postId int) ([]models.Comment, error)
	GetCommentById(commentId int) (models.Comment, error)
	CreateComment(comment models.Comment) error
	DeleteComment(comment models.Comment) error
	ChangeComment(comment models.Comment) error
	// GetPostByCommentId(commentId int) (models.Post, error)
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
	comments, err := s.storage.GetComments(postId)
	if err != nil {
		return nil, fmt.Errorf("service: get commnets: %w", err)
	}
	return comments, nil
}

func (s *CommentService) GetCommentById(commentId int) (models.Comment, error) {
	comment, err := s.storage.GetCommentById(commentId)
	if err != nil {
		return comment, fmt.Errorf("service: get comment by id: %w", err)
	}
	return comment, nil
}

func (s *CommentService) CreateComment(comment models.Comment) error {
	if strings.ReplaceAll(comment.Text, " ", "") == "" {
		return fmt.Errorf("service: create comment: %w", ErrInvalidComment)
	}
	if err := s.storage.CreateComment(comment); err != nil {
		return fmt.Errorf("service: create comment: %w", err)
	}
	return nil
}

func (s *CommentService) DeleteComment(comment models.Comment) error {
	if err := s.storage.DeleteComment(comment); err != nil {
		return fmt.Errorf("service: delete comment: %w", err)
	}
	return nil
}

func (s *CommentService) ChangeComment(comment models.Comment) error {
	if strings.ReplaceAll(comment.Text, " ", "") == "" {
		return fmt.Errorf("service: create comment: %w", ErrInvalidComment)
	}
	if err := s.storage.ChangeComment(comment); err != nil {
		return fmt.Errorf("service: delete comment: %w", err)
	}
	return nil
}

// func (s *CommentService) GetPostByCommentId(commentId int) (models.Post, error) {
// 	return s.storage.GetPostByCommentId(commentId)
// }
