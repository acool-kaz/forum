package storage

import (
	"database/sql"
	"fmt"
)

type LikeDislikeComment interface {
	GetCommentLikes(postId int) (map[int][]string, error)
	LikeComment(commentId int, username string) error
	CommentHasLike(commentId int, username string) error
	RemoveLikeFromComment(commentId int, username string) error
	GetCommentDislikes(postId int) (map[int][]string, error)
	DislikeComment(commentId int, username string) error
	CommentHasDislike(commentId int, username string) error
	RemoveDislikeFromComment(commentId int, username string) error
}

type LikeDislikeCommentStorage struct {
	db *sql.DB
}

func newLikeDislikeCommentStorage(db *sql.DB) *LikeDislikeCommentStorage {
	return &LikeDislikeCommentStorage{
		db: db,
	}
}

func (s *LikeDislikeCommentStorage) GetCommentLikes(postId int) (map[int][]string, error) {
	queryForCommentsId := `SELECT id FROM comment WHERE postId = $1;`
	queryForUsers := `SELECT username FROM likes WHERE commentId = $1;`
	users := make(map[int][]string)
	rowsComment, err := s.db.Query(queryForCommentsId, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get comment likes: %w", err)
	}
	for rowsComment.Next() {
		var id int
		if err := rowsComment.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("storage: get comment likes: %w", err)
		}
		var usernames []string
		rowsUsers, err := s.db.Query(queryForUsers, id)
		if err != nil {
			return nil, fmt.Errorf("storage: get comment likes: %w", err)
		}
		for rowsUsers.Next() {
			var username string
			if err := rowsUsers.Scan(&username); err != nil {
				return nil, fmt.Errorf("storage: get comment likes: %w", err)
			}
			usernames = append(usernames, username)
		}
		users[id] = usernames
	}
	return users, nil
}

func (s *LikeDislikeCommentStorage) LikeComment(commentId int, username string) error {
	query := `INSERT INTO likes(commentId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, commentId, username)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	query = `UPDATE comment SET likes = likes + 1  WHERE id = $1;`
	_, err = s.db.Exec(query, commentId)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	return nil
}

func (s *LikeDislikeCommentStorage) CommentHasLike(commentId int, username string) error {
	var u, query string
	query = `SELECT username FROM likes WHERE commentId = $1 AND username = $2;`
	err := s.db.QueryRow(query, commentId, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: comment has like: %w", err)
	}
	return nil
}

func (s *LikeDislikeCommentStorage) RemoveLikeFromComment(commentId int, username string) error {
	query := `DELETE FROM likes WHERE commentId = $1 AND username = $2;`
	_, err := s.db.Exec(query, commentId, username)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	query = `UPDATE comment SET likes = likes - 1 WHERE id = $1;`
	_, err = s.db.Exec(query, commentId)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	return nil
}

func (s *LikeDislikeCommentStorage) GetCommentDislikes(postId int) (map[int][]string, error) {
	queryForCommentsId := `SELECT id FROM comment WHERE postId = $1;`
	queryForUsers := `SELECT username FROM dislikes WHERE commentId = $1;`
	users := make(map[int][]string)
	rowsComment, err := s.db.Query(queryForCommentsId, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get comment dislikes: %w", err)
	}
	for rowsComment.Next() {
		var id int
		if err := rowsComment.Scan(&id); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("storage: get comment dislikes: %w", err)
		}
		var usernames []string
		rowsUsers, err := s.db.Query(queryForUsers, id)
		if err != nil {
			return nil, fmt.Errorf("storage: get comment dislikes: %w", err)
		}
		for rowsUsers.Next() {
			var username string
			if err := rowsUsers.Scan(&username); err != nil {
				return nil, fmt.Errorf("storage: get comment dislikes: %w", err)
			}
			usernames = append(usernames, username)
		}
		users[id] = usernames
	}
	return users, nil
}

func (s *LikeDislikeCommentStorage) DislikeComment(commentId int, username string) error {
	query := `INSERT INTO dislikes(commentId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, commentId, username)
	if err != nil {
		return fmt.Errorf("storage: dislike comment: %w", err)
	}
	query = `UPDATE comment SET dislikes = dislikes + 1 WHERE id = $1;`
	_, err = s.db.Exec(query, commentId)
	if err != nil {
		return fmt.Errorf("storage: dislike comment: %w", err)
	}
	return nil
}

func (s *LikeDislikeCommentStorage) CommentHasDislike(commentId int, username string) error {
	var u, query string
	query = `SELECT username FROM dislikes WHERE commentId = $1 AND username = $2;`
	err := s.db.QueryRow(query, commentId, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: comment has dislike: %w", err)
	}
	return nil
}

func (s *LikeDislikeCommentStorage) RemoveDislikeFromComment(commentId int, username string) error {
	query := `DELETE FROM dislikes WHERE commentId = $1 AND username = $2;`
	_, err := s.db.Exec(query, commentId, username)
	if err != nil {
		return fmt.Errorf("storage: remove dislike from comment: %w", err)
	}
	query = `UPDATE comment SET dislikes = dislikes - 1  WHERE id = $1;`
	_, err = s.db.Exec(query, commentId)
	if err != nil {
		return fmt.Errorf("storage: remove dislike from comment: %w", err)
	}
	return nil
}
