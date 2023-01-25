package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/config"
	"forum/internal/models"
	"forum/internal/storage"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/google/uuid"
)

type PostService struct {
	postStorage    storage.Post
	tagsStorage    storage.Tags
	commentStorage storage.Comment
	cfg            *config.Config
}

func newPostService(postStorage storage.Post, tagsStorage storage.Tags, commentStorage storage.Comment, cfg *config.Config) *PostService {
	return &PostService{
		postStorage:    postStorage,
		tagsStorage:    tagsStorage,
		commentStorage: commentStorage,
		cfg:            cfg,
	}
}

func (s *PostService) Create(ctx context.Context, post models.Post, files []*multipart.FileHeader) (uint, error) {
	id, err := s.postStorage.Create(ctx, post)
	if err != nil {
		return 0, fmt.Errorf("post service: create: %w", err)
	}

	for _, file := range files {
		path, err := s.saveImages(ctx, id, file)
		if err != nil {
			if err = s.postStorage.Delete(ctx, id); err != nil {
				return 0, fmt.Errorf("post service: create: %w", err)
			}
			return 0, fmt.Errorf("post service: create: %w", err)
		}

		if err = s.postStorage.SaveImages(ctx, id, path); err != nil {
			fmt.Println(err)
			if err = s.postStorage.Delete(ctx, id); err != nil {
				return 0, fmt.Errorf("post service: create: %w", err)
			}
			return 0, fmt.Errorf("post service: create: %w", err)
		}
	}

	for _, name := range strings.Split(post.Tags, " ") {
		if err = s.tagsStorage.Create(ctx, models.Tag{PostId: id, Name: name}); err != nil {
			if err = s.postStorage.Delete(ctx, id); err != nil {
				return 0, fmt.Errorf("post service: create: %w", err)
			}
			return 0, fmt.Errorf("post service: create: %w", err)
		}
	}
	return id, nil
}

func (s *PostService) saveImages(ctx context.Context, postId uint, fileHeader *multipart.FileHeader) (string, error) {
	path := fmt.Sprintf("./static/img/%d", postId)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("post service: save images: %w", err)
	}

	if !strings.Contains(fileHeader.Header["Content-Type"][0], "image") {
		return "", fmt.Errorf("post service: save images: %w", models.ErrInvalidImage)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("post service: save images: %w", err)
	}
	defer file.Close()

	temp := strings.Split(fileHeader.Filename, ".")
	fileType := temp[len(temp)-1]

	if s.isInvalidImageType(fileType) {
		return "", fmt.Errorf("post service: save images: %w", models.ErrInvalidImage)
	}

	fileName := uuid.NewString()
	path += "/" + fileName + "." + fileType

	out, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("post service: save images: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("post service: save images: %w", err)
	}

	return s.cfg.App.FileServer + path[1:], nil
}

func (s *PostService) isInvalidImageType(imageType string) bool {
	validImageType := []string{"jpeg", "jpg", "png"}
	for _, t := range validImageType {
		if t == imageType {
			return false
		}
	}
	return true
}

func (s *PostService) GetAll(ctx context.Context) ([]models.FullPost, error) {
	allPosts, err := s.postStorage.GetAll(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post service: get all: %w", err)
		}
	}
	return allPosts, nil
}

func (s *PostService) GetById(ctx context.Context, id uint) (models.FullPost, error) {
	post, err := s.postStorage.GetById(ctx, id)
	if err != nil {
		return models.FullPost{}, fmt.Errorf("post service: get by id: %w", err)
	}
	post.Comments, err = s.commentStorage.GetAll(ctx, post.Id)
	if err != nil {
		return models.FullPost{}, fmt.Errorf("post service: get by id: %w", err)
	}
	return post, nil
}
