package storage

import (
	"database/sql"
	"fmt"
	"forum/models"
)

type User interface {
	GetPostByUsername(username string) ([]models.Post, error)
	GetLikedPostByUsername(username string) ([]models.Post, error)
	GetCommentedPostByUsername(username string) ([]models.Post, error)
	GetAllCategoryByPostId(postId int) ([]string, error)
	GetUserByUsername(username string) (models.User, error)
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
		return nil, fmt.Errorf("storage: get post by username: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("storage: get post by username: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *UserStorage) GetLikedPostByUsername(username string) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT * FROM post WHERE id IN (SELECT postId FROM likes WHERE username = $1);`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("storage: get like post by username: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("storage: get like post by username: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *UserStorage) GetCommentedPostByUsername(username string) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT * FROM post WHERE id IN (SELECT postId FROM comment WHERE creater = $1);`
	rows, err := s.db.Query(query, username)
	if err != nil {
		return nil, fmt.Errorf("storage: get commented post by username: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes); err != nil {
			return nil, fmt.Errorf("storage: get commented post by username: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *UserStorage) GetAllCategoryByPostId(postId int) ([]string, error) {
	queryCategory := `SELECT category FROM post_category where postId=$1;`
	categoryRows, err := s.db.Query(queryCategory, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get all category by post id: %w", err)
	}
	var category []string
	for categoryRows.Next() {
		var oneCategory string
		if err := categoryRows.Scan(&oneCategory); err != nil {
			return nil, fmt.Errorf("storage: get all category by post id: %w", err)
		}
		category = append(category, oneCategory)
	}
	return category, nil
}

func (s *UserStorage) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	query := `SELECT id, email, username, posts, likes, comments FROM user WHERE username = $1;`
	row := s.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.CountOfPosts, &user.CountOfLikes, &user.CountOfComments)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by username: %w", err)
	}
	return user, err
}
