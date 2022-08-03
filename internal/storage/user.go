package storage

import (
	"database/sql"
	"forum/models"
)

type User interface {
	GetPostByUsername(username string) ([]models.Post, error)
	GetAllCategoryByPostId(postId int) ([]string, error)
}

type UserStorage struct {
	db *sql.DB
}

func newUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) GetPostByUsername(username string) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT * FROM post WHERE creater = $1;`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *UserStorage) GetAllCategoryByPostId(postId int) ([]string, error) {
	queryCategory := `SELECT category FROM post_category where postId=$1;`
	categoryRows, err := s.db.Query(queryCategory, postId)
	if err != nil {
		return nil, err
	}
	var category []string
	for categoryRows.Next() {
		var oneCategory string
		if err := categoryRows.Scan(&oneCategory); err != nil {
			return nil, err
		}
		category = append(category, oneCategory)
	}
	return category, nil
}
