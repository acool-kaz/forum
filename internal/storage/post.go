package storage

import (
	"context"
	"database/sql"
	"fmt"
	"forum/models"
	"strings"
)

type PostStorage struct {
	db *sql.DB
}

func newPostStorage(db *sql.DB) *PostStorage {
	return &PostStorage{
		db: db,
	}
}

func (s *PostStorage) Create(ctx context.Context, post models.Post) (uint, error) {
	query := fmt.Sprintf("INSERT INTO %s(user_id, title, description) VALUES ($1, $2, $3) RETURNING id;", postTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("post storage: create: %w", err)
	}
	defer prep.Close()

	var id uint
	if err = prep.QueryRowContext(ctx, post.UserId, post.Title, post.Description).Scan(&id); err != nil {
		return 0, fmt.Errorf("post storage: crete: %w", err)
	}

	return id, nil
}

func (s *PostStorage) GetAll(ctx context.Context) ([]models.FullPost, error) {
	query := fmt.Sprintf(`
	SELECT 
		p.id,
		u.username,
		p.title,
		GROUP_CONCAT(t.name, ' '),
		p.description,
		p.created_at
	FROM %s p 
	INNER JOIN %s u ON u.id = p.user_id
	INNER JOIN %s t ON t.post_id = p.id
	GROUP BY p.id;
	`, postTable, userTable, tagTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("post storage: get all: %w", err)
	}
	defer prep.Close()

	var (
		allPosts []models.FullPost
		onePost  models.FullPost
		tags     string
	)

	rows, err := prep.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("post storage: get all: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&onePost.Id, &onePost.Username, &onePost.Title, &tags, &onePost.Description, &onePost.CreatedAt); err != nil {
			return nil, fmt.Errorf("post storage: get all: %w", err)
		}
		onePost.Tags = strings.Split(tags, " ")
		allPosts = append(allPosts, onePost)
	}

	return allPosts, nil
}

func (s *PostStorage) GetById(ctx context.Context, id uint) (models.FullPost, error) {
	query := fmt.Sprintf(`
	SELECT 
		p.id,
		u.username,
		p.title,
		GROUP_CONCAT(t.name, ' '),
		p.description,
		(SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND reaction=1) AS 'likes',
    	(SELECT COUNT(*) FROM reactions r WHERE r.post_id = p.id AND reaction=-1) AS 'dislikes',
		p.created_at
	FROM %s p 
	INNER JOIN %s u ON u.id = p.user_id
	INNER JOIN %s t ON t.post_id = p.id
	WHERE p.id = $1
	GROUP BY p.id;
	`, postTable, userTable, tagTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return models.FullPost{}, fmt.Errorf("post storage: get all: %w", err)
	}
	defer prep.Close()

	var (
		onePost models.FullPost
		tags    string
	)

	if err = prep.QueryRowContext(ctx, id).Scan(&onePost.Id, &onePost.Username, &onePost.Title, &tags, &onePost.Description, &onePost.Likes, &onePost.Dislikes, &onePost.CreatedAt); err != nil {
		return models.FullPost{}, fmt.Errorf("post storage: get all: %w", err)
	}
	onePost.Tags = strings.Split(tags, " ")

	return onePost, nil
}

func (s *PostStorage) Delete(ctx context.Context, id uint) error {
	query := fmt.Sprintf("DELTE FROM %s WHERE id = $1;", postTable)
	prep, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("post storage: delete: %w", err)
	}
	defer prep.Close()

	if _, err = prep.ExecContext(ctx, id); err != nil {
		return fmt.Errorf("post storage: delete: %w", err)
	}

	return nil
}
