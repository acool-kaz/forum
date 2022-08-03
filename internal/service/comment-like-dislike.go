package service

import "forum/internal/storage"

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
		return s.storage.RemoveLikeFromComment(commentId, username)
	}
	if err := s.storage.CommentHasDislike(commentId, username); err == nil {
		s.storage.RemoveDislikeFromComment(commentId, username)
	}
	return s.storage.LikeComment(commentId, username)
}

func (s *LikeDislikeCommentService) GetCommentLikes(postId int) (map[int][]string, error) {
	return s.storage.GetCommentLikes(postId)
}

func (s *LikeDislikeCommentService) DislikeComment(commentId int, username string) error {
	if err := s.storage.CommentHasDislike(commentId, username); err == nil {
		return s.storage.RemoveDislikeFromComment(commentId, username)
	}
	if err := s.storage.CommentHasLike(commentId, username); err == nil {
		s.storage.RemoveLikeFromComment(commentId, username)
	}
	return s.storage.DislikeComment(commentId, username)
}

func (s *LikeDislikeCommentService) GetCommentDislikes(postId int) (map[int][]string, error) {
	return s.storage.GetCommentDislikes(postId)
}
