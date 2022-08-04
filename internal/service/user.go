package service

import (
	"forum/internal/storage"
	"forum/models"
)

type User interface {
	GetPostByUsername(username string) ([]models.Post, error)
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

func (s *UserService) GetPostByUsername(username string) ([]models.Post, error) {
	posts, err := s.storage.GetPostByUsername(username)
	if err != nil {
		return nil, err
	}
	for i := range posts {
		category, err := s.storage.GetAllCategoryByPostId(posts[i].Id)
		if err != nil {
			return nil, err
		}
		posts[i].Category = category
	}
	return posts, nil
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
	return s.storage.GetUserByUsername(username)
}
