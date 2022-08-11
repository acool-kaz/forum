package service

import (
	"errors"
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"strings"
)

var ErrInvalidQuery = errors.New("invalid query request")

type User interface {
	GetPostByUsername(username string, query map[string][]string) ([]models.Post, error)
	GetUserByUsername(username string) (models.User, error)
}

type UserService struct {
	storage storage.User
}

func newUserService(storage storage.User) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) GetPostByUsername(username string, query map[string][]string) ([]models.Post, error) {
	var (
		posts []models.Post
		err   error
	)
	search, ok := query["posts"]
	if !ok {
		return nil, fmt.Errorf("service: get post by username: %w", ErrInvalidQuery)
	}
	switch strings.Join(search, "") {
	case "created":
		posts, err = s.storage.GetPostByUsername(username)
	case "liked":
		posts, err = s.storage.GetLikedPostByUsername(username)
	case "commented":
		posts, err = s.storage.GetCommentedPostByUsername(username)
	default:
		return nil, fmt.Errorf("service: get post by username: %w", ErrInvalidQuery)
	}
	if err != nil {
		return nil, fmt.Errorf("service: get post by username: %w", err)
	}
	for i := range posts {
		category, err := s.storage.GetAllCategoryByPostId(posts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("service: get post by username: %w", err)
		}
		posts[i].Category = category
	}
	return posts, nil
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	return s.storage.GetUserByUsername(username)
}
