package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/storage"
)

type LikeDislikeComment interface {
	LikeComment(commentId int, username string) error
	GetCommentLikes(postId int) (map[int][]string, error)
	DislikeComment(commentId int, username string) error
	GetCommentDislikes(postId int) (map[int][]string, error)
}

type LikeDislikeCommentService struct {
	storage storage.LikeDislikeComment
}

func newLikeDislikeCommentService(storage storage.LikeDislikeComment) *LikeDislikeCommentService {
	return &LikeDislikeCommentService{
		storage: storage,
	}
}

func (s *LikeDislikeCommentService) LikeComment(commentId int, username string) error {
	if err := s.storage.CommentHasLike(commentId, username); err == nil {
		if err := s.storage.RemoveLikeFromComment(commentId, username); err != nil {
			return fmt.Errorf("service: like comment: %w", err)
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like comment: %w", err)
	}
	if err := s.storage.CommentHasDislike(commentId, username); err == nil {
		if err := s.storage.RemoveDislikeFromComment(commentId, username); err != nil {
			return fmt.Errorf("service: like comment: %w", err)
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like comment: %w", err)
	}
	if err := s.storage.LikeComment(commentId, username); err != nil {
		return fmt.Errorf("service: like comment: %w", err)
	}
	return nil
}

func (s *LikeDislikeCommentService) GetCommentLikes(postId int) (map[int][]string, error) {
	users, err := s.storage.GetCommentLikes(postId)
	if err != nil {
		return nil, fmt.Errorf("service: get comment likes: %w", err)
	}
	return users, nil
}

func (s *LikeDislikeCommentService) DislikeComment(commentId int, username string) error {
	if err := s.storage.CommentHasDislike(commentId, username); err == nil {
		if err := s.storage.RemoveDislikeFromComment(commentId, username); err != nil {
			return fmt.Errorf("service: dislike comment: %w", err)
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike comment: %w", err)
	}
	if err := s.storage.CommentHasLike(commentId, username); err == nil {
		if err := s.storage.RemoveLikeFromComment(commentId, username); err != nil {
			return fmt.Errorf("service: dislike comment: %w", err)
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike comment: %w", err)
	}
	if err := s.storage.DislikeComment(commentId, username); err != nil {
		return fmt.Errorf("service: dislike comment: %w", err)
	}
	return nil
}

func (s *LikeDislikeCommentService) GetCommentDislikes(postId int) (map[int][]string, error) {
	users, err := s.storage.GetCommentDislikes(postId)
	if err != nil {
		return nil, fmt.Errorf("service: get comment likes: %w", err)
	}
	return users, nil
}
