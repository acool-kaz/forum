package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/config"
	"forum/internal/models"
	"forum/internal/storage"
	"forum/pkg/image_saver"
	"mime/multipart"
	"strconv"
	"strings"
)

type PostService struct {
	postStorage    storage.Post
	tagsStorage    storage.Tags
	commentStorage storage.Comment
	imageStorage   storage.Image
	cfg            *config.Config
}

func newPostService(postStorage storage.Post, tagsStorage storage.Tags, commentStorage storage.Comment, imageStorage storage.Image, cfg *config.Config) *PostService {
	return &PostService{
		postStorage:    postStorage,
		tagsStorage:    tagsStorage,
		commentStorage: commentStorage,
		imageStorage:   imageStorage,
		cfg:            cfg,
	}
}

func (s *PostService) Update(ctx context.Context, userId, postId uint, updateInfo models.UpdatePost) error {
	// TODO:
	// get post for user
	// check if it his post
	return nil
}

func (s *PostService) Create(ctx context.Context, post models.Post, files []*multipart.FileHeader) (uint, error) {
	id, err := s.postStorage.Create(ctx, post)
	if err != nil {
		return 0, fmt.Errorf("post service: create: %w", err)
	}

	for _, file := range files {
		path, err := image_saver.SaveImages(ctx, "./static/img/", strconv.Itoa(int(id)), file)
		if err != nil {
			if err = s.postStorage.Delete(ctx, id); err != nil {
				return 0, fmt.Errorf("post service: create: %w", err)
			}
			return 0, fmt.Errorf("post service: create: %w", err)
		}

		if err = s.imageStorage.SaveImages(ctx, id, s.cfg.App.FileServer+path); err != nil {
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

func (s *PostService) GetAll(ctx context.Context) ([]models.FullPost, error) {
	allPosts, err := s.postStorage.GetAll(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post service: get all: %w", err)
		}
	}

	profilePostFilter := ctx.Value(models.ProfilePostFilter)
	if profilePostFilter != nil && profilePostFilter.(string) == "commented" {
		for i, post := range allPosts {
			allPosts[i].Comments, err = s.commentStorage.GetAll(ctx, post.Id)
			if err != nil {
				return nil, fmt.Errorf("post service: get all: %w", err)
			}
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
