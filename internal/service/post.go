package service

import (
	"errors"
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"strings"
)

var (
	ErrInvalidPost         = errors.New("invalid post")
	ErrInvalidQueryRequest = errors.New("invalid query request")
)

type Post interface {
	CreatePost(post models.Post) (int, error)
	GetAllPost() ([]models.Post, error)
	GetPostById(postId int) (models.Post, error)
	GetAllPostBy(user models.User, query map[string][]string) ([]models.Post, error)
	GetSimilarPosts(postId int) ([]models.Post, error)
	DeletePost(post models.Post) error
	ChangePost(newPost models.Post, postId int) error
	SaveImageForPost(postId int, imgPath string) error
}

type PostService struct {
	storage storage.Post
}

func newPostService(storage storage.Post) *PostService {
	return &PostService{
		storage: storage,
	}
}

func (s *PostService) CreatePost(post models.Post) (int, error) {
	post.Category = strings.Fields(strings.Join(append(post.Category[1:], post.Category[:1]...), " "))
	post.Category = deleteDuplicate(post.Category)
	if isInvalidPost(post) {
		return 0, fmt.Errorf("service: create post: %w", ErrInvalidPost)
	}
	id, err := s.storage.CreatePost(post)
	if err != nil {
		return 0, fmt.Errorf("service: create post: %w", err)
	}
	return id, nil
}

func (s *PostService) GetAllPost() ([]models.Post, error) {
	posts, err := s.storage.GetAllPost()
	if err != nil {
		return nil, fmt.Errorf("service: get all post: %w", err)
	}
	for i := range posts {
		category, err := s.storage.GetAllCategoryByPostId(posts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("service: get all post: %w", err)
		}
		posts[i].Category = category
	}
	return posts, nil
}

func (s *PostService) GetPostById(postId int) (models.Post, error) {
	post, err := s.storage.GetPostById(postId)
	if err != nil {
		return post, fmt.Errorf("service: get post by id: %w", err)
	}
	post.Category, err = s.storage.GetAllCategoryByPostId(post.Id)
	if err != nil {
		return post, fmt.Errorf("service: get post by id: %w", err)
	}
	post.Images, err = s.storage.GetAllImagesByPostId(post.Id)
	if err != nil {
		return post, fmt.Errorf("service: get post by id: %w", err)
	}
	return post, nil
}

func (s *PostService) GetAllPostBy(user models.User, query map[string][]string) ([]models.Post, error) {
	var (
		posts []models.Post
		err   error
	)
	for key, val := range query {
		switch key {
		case "category":
			posts, err = s.storage.GetPostsByCategory(strings.Join(val, ""))
			if err != nil {
				return nil, fmt.Errorf("service: get all post by: %w", err)
			}
		case "time":
			switch strings.Join(val, "") {
			case "new":
				posts, err = s.storage.GetPostByTimeNew()
			case "old":
				posts, err = s.storage.GetPostByTimeOld()
			default:
				return nil, fmt.Errorf("service: get all post by: %w", ErrInvalidQueryRequest)
			}
			if err != nil {
				return nil, fmt.Errorf("service: get all post by: %w", err)
			}
		case "likes":
			switch strings.Join(val, "") {
			case "most":
				posts, err = s.storage.GetPostByLikeMost()
			case "least":
				posts, err = s.storage.GetPostByLikeLeast()
			default:
				return nil, fmt.Errorf("service: get all post by: %w", ErrInvalidQueryRequest)
			}
			if err != nil {
				return nil, fmt.Errorf("service: get all post by: %w", err)
			}
		default:
			return nil, fmt.Errorf("service: get all post by: %w", ErrInvalidQueryRequest)
		}
	}
	for i := range posts {
		category, err := s.storage.GetAllCategoryByPostId(posts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("service: get all post by: %w", err)
		}
		posts[i].Category = category
	}
	return posts, nil
}

func (s *PostService) GetSimilarPosts(postId int) ([]models.Post, error) {
	posts, err := s.storage.GetSimilarPosts(postId)
	if err != nil {
		return nil, fmt.Errorf("service: get similar posts: %w", err)
	}
	return posts, nil
}

func (s *PostService) DeletePost(post models.Post) error {
	if err := s.storage.DeletePost(post); err != nil {
		return fmt.Errorf("service: delete post: %w", err)
	}
	return nil
}

func (s *PostService) ChangePost(newPost models.Post, postId int) error {
	if isInvalidPost(newPost) {
		return fmt.Errorf("service: create post: %w", ErrInvalidPost)
	}
	if err := s.storage.ChangePost(newPost, postId); err != nil {
		return fmt.Errorf("service: change post: %w", err)
	}
	return nil
}

func (s *PostService) SaveImageForPost(postId int, imgPath string) error {
	if err := s.storage.SaveImageForPost(postId, imgPath); err != nil {
		return fmt.Errorf("service: save image for post: %w", err)
	}
	return nil
}

func isInvalidPost(post models.Post) bool {
	if strings.ReplaceAll(post.Title, " ", "") == "" {
		return true
	}
	if strings.ReplaceAll(post.Description, " ", "") == "" {
		return true
	}
	if len(post.Category) == 0 {
		return true
	}
	return false
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
