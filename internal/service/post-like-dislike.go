package service

import "forum/internal/storage"

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
		return s.storage.RemoveLikeFromPost(postId, username)
	}
	if err := s.storage.PostHasDislike(postId, username); err == nil {
		s.storage.RemoveDislikeFromPost(postId, username)
	}
	return s.storage.LikePost(postId, username)
}

func (s *ads) GetPostLikes(postId int) ([]string, error) {
	return s.storage.GetPostLikes(postId)
}

func (s *ads) DislikePost(postId int, username string) error {
	if err := s.storage.PostHasDislike(postId, username); err == nil {
		return s.storage.RemoveDislikeFromPost(postId, username)
	}
	if err := s.storage.PostHasLike(postId, username); err == nil {
		s.storage.RemoveLikeFromPost(postId, username)
	}
	return s.storage.DislikePost(postId, username)
}

func (s *ads) GetPostDislikes(postId int) ([]string, error) {
	return s.storage.GetPostDislikes(postId)
}
