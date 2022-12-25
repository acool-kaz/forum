package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/models"
)

type TagsStorage struct {
	db *sql.DB
}

func newTagsStorage(db *sql.DB) *TagsStorage {
	return &TagsStorage{
		db: db,
	}
}

func (s *TagsStorage) Create(ctx context.Context, tag models.Tag) error {
	query := fmt.Sprintf("INSERT INTO %s(post_id, name) VALUES ($1, $2);", tagTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("tags storage: create: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, tag.PostId, tag.Name); err != nil {
		return fmt.Errorf("tags storage: create: %w", err)
	}

	return nil
}
