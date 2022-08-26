package storage

import (
	"database/sql"
	"fmt"
)

type LikeDislikePost interface {
	GetPostLikes(postId int) ([]string, error)
	LikePost(postId int, username string) error
	PostHasLike(postId int, username string) error
	RemoveLikeFromPost(postId int, username string) error
	GetPostDislikes(postId int) ([]string, error)
	DislikePost(postId int, username string) error
	PostHasDislike(postId int, username string) error
	RemoveDislikeFromPost(postId int, username string) error
}

type LikeDislikePostStorage struct {
	db *sql.DB
}

func newLikeDislikePostStorage(db *sql.DB) *LikeDislikePostStorage {
	return &LikeDislikePostStorage{
		db: db,
	}
}

func (s *LikeDislikePostStorage) GetPostLikes(postId int) ([]string, error) {
	var postLikes []string
	query := `SELECT username FROM likes WHERE postId = $1;`
	rows, err := s.db.Query(query, postId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var postLike string
		if err := rows.Scan(&postLike); err != nil {
			if err == sql.ErrNoRows {
				return []string{}, nil
			}
			return nil, err
		}
		postLikes = append(postLikes, postLike)
	}
	return postLikes, nil
}

func (s *LikeDislikePostStorage) LikePost(postId int, username string) error {
	query := `INSERT INTO likes(postId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, postId, username)
	if err != nil {
		return fmt.Errorf("storage: like post: %w", err)
	}
	query = `UPDATE post SET likes = likes + 1 WHERE id = $1;`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: like post: %w", err)
	}
	query = `UPDATE user SET likes = likes + 1 WHERE username = (SELECT creater FROM post WHERE id = $1);`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: like post: %w", err)
	}
	return err
}

func (s *LikeDislikePostStorage) PostHasLike(postId int, username string) error {
	var u, query string
	query = `SELECT username FROM likes WHERE postId = $1 AND username = $2;`
	err := s.db.QueryRow(query, postId, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: post has like: %w", err)
	}
	return nil
}

func (s *LikeDislikePostStorage) RemoveLikeFromPost(postId int, username string) error {
	query := `DELETE FROM likes WHERE postId = $1 AND username = $2;`
	_, err := s.db.Exec(query, postId, username)
	if err != nil {
		return fmt.Errorf("storage: remove like from post: %w", err)
	}
	query = `UPDATE post SET likes = likes - 1 WHERE id = $1;`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: remove like from post: %w", err)
	}
	query = `UPDATE user SET likes = likes - 1 WHERE username = (SELECT creater FROM post WHERE id = $1);`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: remove like from post: %w", err)
	}
	return nil
}

func (s *LikeDislikePostStorage) GetPostDislikes(postId int) ([]string, error) {
	var postDislikes []string
	query := `SELECT username FROM dislikes WHERE postId = $1;`
	rows, err := s.db.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get post dislikes: %w", err)
	}
	for rows.Next() {
		var postDislike string
		if err := rows.Scan(&postDislike); err != nil {
			if err == sql.ErrNoRows {
				return []string{}, nil
			}
			return nil, fmt.Errorf("storage: get post dislikes: %w", err)
		}
		postDislikes = append(postDislikes, postDislike)
	}
	return postDislikes, nil
}

func (s *LikeDislikePostStorage) DislikePost(postId int, username string) error {
	query := `INSERT INTO dislikes(postId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, postId, username)
	if err != nil {
		return fmt.Errorf("storage: dislike post: %w", err)
	}
	query = `UPDATE post SET dislikes = dislikes + 1 WHERE id = $1;`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: dislike post: %w", err)
	}
	query = `UPDATE user SET dislikes = dislikes + 1 WHERE username = (SELECT creater FROM post WHERE id = $1);`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: remove dislike from post: %w", err)
	}
	return nil
}

func (s *LikeDislikePostStorage) PostHasDislike(postId int, username string) error {
	var u, query string
	query = `SELECT username FROM dislikes WHERE postId = $1 AND username = $2;`
	err := s.db.QueryRow(query, postId, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: post has dislike: %w", err)
	}
	return nil
}

func (s *LikeDislikePostStorage) RemoveDislikeFromPost(postId int, username string) error {
	query := `DELETE FROM dislikes WHERE postId = $1 AND username = $2;`
	_, err := s.db.Exec(query, postId, username)
	if err != nil {
		return fmt.Errorf("storage: remove dislike from post: %w", err)
	}
	query = `UPDATE post SET dislikes = dislikes - 1 WHERE id = $1;`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: remove dislike from post: %w", err)
	}
	query = `UPDATE user SET dislikes = dislikes - 1 WHERE username = (SELECT creater FROM post WHERE id = $1);`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: remove dislike from post: %w", err)
	}
	return nil
}
