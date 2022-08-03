package service

import (
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"strings"
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPost() ([]models.Post, error)
	GetPostById(postId int) (models.Post, error)
	GetPostsByCategory(category string) ([]models.Post, error)
	GetPostsByTime(time string) ([]models.Post, error)
	GetPostsByLike(like string) ([]models.Post, error)
	GetSimilarPosts(postId int) ([]models.Post, error)
}

type PostService struct {
	storage storage.Post
}

func newPostService(storage storage.Post) *PostService {
	return &PostService{
		storage: storage,
	}
}

func (s *PostService) GetAllPost() ([]models.Post, error) {
	posts, err := s.storage.GetAllPost()
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

func (s *PostService) GetPostsByCategory(category string) ([]models.Post, error) {
	posts, err := s.storage.GetPostsByCategory(category)
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

func (s *PostService) GetPostsByTime(time string) ([]models.Post, error) {
	var posts []models.Post
	var err error
	switch time {
	case "new":
		posts, err = s.storage.GetPostByTimeNew()
		if err != nil {
			return nil, err
		}
	case "old":
		posts, err = s.storage.GetPostByTimeOld()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid argument for filter by time")
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

func (s *PostService) GetPostsByLike(like string) ([]models.Post, error) {
	var posts []models.Post
	var err error
	switch like {
	case "most":
		posts, err = s.storage.GetPostByLikeMost()
		if err != nil {
			return nil, err
		}
	case "least":
		posts, err = s.storage.GetPostByLikeLeast()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid argument for filter by time")
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

func (s *PostService) GetSimilarPosts(postId int) ([]models.Post, error) {
	return s.storage.GetSimilarPosts(postId)
}

func (s *PostService) CreatePost(post models.Post) error {
	post.Description = strings.TrimPrefix(post.Description, "\r\n")
	post.Category = strings.Fields(strings.Join(append(post.Category[1:], post.Category[:1]...), " "))
	post.Category = deleteDuplicate(post.Category)
	return s.storage.CreatePost(post)
}

func (s *PostService) GetPostById(postId int) (models.Post, error) {
	post, err := s.storage.GetPostById(postId)
	if err != nil {
		return models.Post{}, err
	}
	post.Category, err = s.storage.GetAllCategoryByPostId(post.Id)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func deleteDuplicate(arr []string) []string {
	temp := []string{}
mainloop:
	for _, item := range arr {
		for _, tempItem := range temp {
			if tempItem == item {
				continue mainloop
			}
		}
		temp = append(temp, item)
	}
	return temp
}
