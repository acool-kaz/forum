package service

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/storage"
)

type LikeDislikePost interface {
	LikePost(postId int, username string) error
	GetPostLikes(postId int) ([]string, error)
	DislikePost(postId int, username string) error
	GetPostDislikes(postId int) ([]string, error)
}

type ads struct {
	storage storage.LikeDislikePost
}

func newLikeDislikeService(storage storage.LikeDislikePost) *ads {
	return &ads{
		storage: storage,
	}
}

func (s *ads) LikePost(postId int, username string) error {
	if err := s.storage.PostHasLike(postId, username); err == nil {
		if err := s.storage.RemoveLikeFromPost(postId, username); err != nil {
			return fmt.Errorf("service: like post: %w", err)
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like post: %w", err)
	}
	if err := s.storage.PostHasDislike(postId, username); err == nil {
		if err := s.storage.RemoveDislikeFromPost(postId, username); err != nil {
			return fmt.Errorf("service: like post: %w", err)
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: like post: %w", err)
	}
	if err := s.storage.LikePost(postId, username); err != nil {
		return fmt.Errorf("service: like post: %w", err)
	}
	return nil
}

func (s *ads) GetPostLikes(postId int) ([]string, error) {
	users, err := s.storage.GetPostLikes(postId)
	if err != nil {
		return nil, fmt.Errorf("service: get post likes: %w", err)
	}
	return users, nil
}

func (s *ads) DislikePost(postId int, username string) error {
	if err := s.storage.PostHasDislike(postId, username); err == nil {
		if err := s.storage.RemoveDislikeFromPost(postId, username); err != nil {
			return fmt.Errorf("service: dislike post: %w", err)
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike post: %w", err)
	}
	if err := s.storage.PostHasLike(postId, username); err == nil {
		if err := s.storage.RemoveLikeFromPost(postId, username); err != nil {
			return fmt.Errorf("service: dislike post: %w", err)
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("service: dislike post: %w", err)
	}
	if err := s.storage.DislikePost(postId, username); err != nil {
		return fmt.Errorf("service: dislike post: %w", err)
	}
	return nil
}

func (s *ads) GetPostDislikes(postId int) ([]string, error) {
	users, err := s.storage.GetPostDislikes(postId)
	if err != nil {
		return nil, fmt.Errorf("service: get post dislikes: %w", err)
	}
	return users, nil
}
