package storage

import (
	"database/sql"
	"fmt"
	"forum/models"
)

type Comment interface {
	GetComments(postId int) ([]models.Comment, error)
	CreateComment(comment models.Comment) error
	GetPostIdByCommentId(commentId int) (int, error)
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
		return nil, err
	}
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.Creater, &comment.Text, &comment.Likes, &comment.Dislikes); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		commnets = append(commnets, comment)
	}
	return commnets, nil
}

func (s *CommentStorage) CreateComment(comment models.Comment) error {
	query := `INSERT INTO comment(postId, creater, text) VALUES ($1, $2, $3);`
	_, err := s.db.Exec(query, comment.PostId, comment.Creater, comment.Text)
	return err
}

func (s *CommentStorage) GetPostIdByCommentId(commentId int) (int, error) {
	var postId int
	query := `SELECT postId FROM comment WHERE id = $1;`
	row := s.db.QueryRow(query, commentId)
	err := row.Scan(&postId)
	return postId, err
}
