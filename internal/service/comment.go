package service

import (
	"context"
	"fmt"
	"forum/internal/models"
	"forum/internal/storage"
)

type CommentService struct {
	commentStorage storage.Comment
}

func newCommentService(commentStorage storage.Comment) *CommentService {
	return &CommentService{
		commentStorage: commentStorage,
	}
}

func (s *CommentService) Create(ctx context.Context, comment models.Comment) error {
	if err := s.commentStorage.Create(ctx, comment); err != nil {
		return fmt.Errorf("comment service: create: %w", err)
	}
	return nil
}
