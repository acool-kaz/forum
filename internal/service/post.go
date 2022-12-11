package service

import (
	"errors"
	"fmt"
	"forum/internal/storage"
	"forum/models"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

var (
	ErrInvalidPost         = errors.New("invalid post")
	ErrInvalidQueryRequest = errors.New("invalid query request")
	ErrInavlidUser         = errors.New("invalid user")
)

type Post interface {
	CreatePost(post models.Post) error
	GetAllPost() ([]models.Post, error)
	GetPostById(postId int) (models.Post, error)
	GetAllPostBy(user models.User, query map[string][]string) ([]models.Post, error)
	GetSimilarPosts(postId int) ([]models.Post, error)
	DeletePost(post models.Post, user models.User) error
	ChangePost(newPost, oldPost models.Post, user models.User) error
	SaveImageForPost(postId int, files []*multipart.FileHeader) error
}

type PostService struct {
	storage storage.Post
}

func newPostService(storage storage.Post) *PostService {
	return &PostService{
		storage: storage,
	}
}

func (s *PostService) CreatePost(post models.Post) error {
	post.Category = strings.Fields(strings.Join(append(post.Category[1:], post.Category[:1]...), " "))
	post.Category = deleteDuplicate(post.Category)
	if isInvalidPost(post) {
		return fmt.Errorf("service: create post: %w", ErrInvalidPost)
	}
	id, err := s.storage.CreatePost(post)
	if err != nil {
		return fmt.Errorf("service: create post: %w", err)
	}
	post.Id = id
	if err := s.SaveImageForPost(post.Id, post.Files); err != nil {
		if err := s.storage.DeletePost(post); err != nil {
			return fmt.Errorf("service: create post: %w", err)
		}
		return fmt.Errorf("service: create post: %w", err)
	}
	return nil
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

func (s *PostService) DeletePost(post models.Post, user models.User) error {
	if user.Username != post.Creater {
		return fmt.Errorf("service: delete post: %w: cant delete post", ErrInavlidUser)
	}
	if err := s.storage.DeletePost(post); err != nil {
		return fmt.Errorf("service: delete post: %w", err)
	}
	return nil
}

func (s *PostService) ChangePost(newPost, oldPost models.Post, user models.User) error {
	if user.Username != oldPost.Creater {
		return fmt.Errorf("service: change post: %w: you cant change post", ErrInavlidUser)
	}
	if isInvalidPost(newPost) {
		return fmt.Errorf("service: change post: %w", ErrInvalidPost)
	}
	if err := s.storage.ChangePost(newPost, oldPost.Id); err != nil {
		return fmt.Errorf("service: change post: %w", err)
	}
	return nil
}

func (s *PostService) SaveImageForPost(postId int, files []*multipart.FileHeader) error {
	if err := os.MkdirAll(fmt.Sprintf("./ui/static/img/%d", postId), os.ModePerm); err != nil {
		return fmt.Errorf("service: save image for post: %w", err)
	}
	for _, fileHeader := range files {
		if !strings.Contains(fileHeader.Header["Content-Type"][0], "image") {
			continue
		}
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("service: save image for post: %w", err)
		}
		defer file.Close()
		imageSplit := strings.Split(fileHeader.Filename, ".")
		if !validImageType(imageSplit[len(imageSplit)-1]) {
			return fmt.Errorf("service: save image for post: %w: not supported file type", ErrInvalidPost)
		}
		out, err := os.Create(fmt.Sprintf("./ui/static/img/%d/%s", postId, fileHeader.Filename))
		if err != nil {
			return fmt.Errorf("service: save image for post: %w", err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			return fmt.Errorf("service: save image for post: %w", err)
		}
		if err := s.storage.SaveImageForPost(postId, fmt.Sprintf("/static/img/%d/%s", postId, fileHeader.Filename)); err != nil {
			return fmt.Errorf("service: save image for post: %w", err)
		}
	}
	return nil
}

func validImageType(imageType string) bool {
	validImageType := []string{"jpeg", "jpg", "png", "gif"}
	for _, t := range validImageType {
		if t == imageType {
			return true
		}
	}
	return false
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
