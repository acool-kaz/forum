package storage

import (
	"database/sql"
	"fmt"
	"forum/models"
)

type Comment interface {
	GetComments(postId int) ([]models.Comment, error)
	GetCommentById(commentId int) (models.Comment, error)
	CreateComment(comment models.Comment) error
}

type CommentStorage struct {
	db *sql.DB
}

func newCommentStorage(db *sql.DB) *CommentStorage {
	return &CommentStorage{
		db: db,
	}
}

func (s *CommentStorage) GetComments(postId int) ([]models.Comment, error) {
	var commnets []models.Comment
	query := `SELECT * FROM comment WHERE postId=$1;`
	rows, err := s.db.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get comments: %w", err)
	}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.Creater, &comment.Text, &comment.Likes, &comment.Dislikes); err != nil {
			return nil, fmt.Errorf("storage: get comments: %w", err)
		}
		commnets = append(commnets, comment)
	}
	return commnets, nil
}

func (s *CommentStorage) GetCommentById(commentId int) (models.Comment, error) {
	var comment models.Comment
	query := `SELECT id, creater, postId FROM comment WHERE id = $1;`
	row := s.db.QueryRow(query, commentId)
	err := row.Scan(&comment.Id, &comment.Creater, &comment.PostId)
	if err != nil {
		return models.Comment{}, fmt.Errorf("storage: get comment by id: %w", err)
	}
	return comment, nil
}

func (s *CommentStorage) CreateComment(comment models.Comment) error {
	query := `INSERT INTO comment(postId, creater, text) VALUES ($1, $2, $3);`
	_, err := s.db.Exec(query, comment.PostId, comment.Creater, comment.Text)
	if err != nil {
		return fmt.Errorf("storage: create comment: %w", err)
	}
	query = `UPDATE user SET comments = comments + 1 WHERE username = (SELECT creater FROM post WHERE id = $1);`
	_, err = s.db.Exec(query, comment.PostId)
	if err != nil {
		return fmt.Errorf("storage: create comment: %w", err)
	}
	query = `UPDATE post SET comments = comments + 1 WHERE id = $1;`
	_, err = s.db.Exec(query, comment.PostId)
	if err != nil {
		return fmt.Errorf("storage: create comment: %w", err)
	}
	return nil
}
