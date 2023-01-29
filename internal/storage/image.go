package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type ImageStorage struct {
	db *sql.DB
}

func newImageStorage(db *sql.DB) *ImageStorage {
	return &ImageStorage{
		db: db,
	}
}

func (s *ImageStorage) SaveImages(ctx context.Context, postId uint, url string) error {
	query := fmt.Sprintf("INSERT INTO %s(post_id, url) VALUES ($1, $2);", imageTable)

	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("image storage: save images: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, postId, url); err != nil {
		return fmt.Errorf("image storage: save images: %w", err)
	}

	return nil
}
