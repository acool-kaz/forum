package storage

import (
	"database/sql"
	"fmt"
	"forum/models"
)

type Post interface {
	CreatePost(post models.Post) (int, error)
	GetAllPost() ([]models.Post, error)
	GetPostById(postId int) (models.Post, error)
	GetPostsByCategory(category string) ([]models.Post, error)
	GetPostByTimeNew() ([]models.Post, error)
	GetPostByTimeOld() ([]models.Post, error)
	GetPostByLikeMost() ([]models.Post, error)
	GetPostByLikeLeast() ([]models.Post, error)
	GetSimilarPosts(postId int) ([]models.Post, error)
	GetAllCategoryByPostId(postId int) ([]string, error)
	DeletePost(post models.Post) error
	ChangePost(newPost models.Post, postId int) error
	GetAllImagesByPostId(postId int) ([]string, error)
	SaveImageForPost(postId int, imgPath string) error
}

type PostStorage struct {
	db *sql.DB
}

func newPostStorage(db *sql.DB) *PostStorage {
	return &PostStorage{
		db: db,
	}
}

func (s *PostStorage) CreatePost(post models.Post) (int, error) {
	query := `INSERT INTO post (creater, title, description) VALUES ($1, $2, $3);`
	result, err := s.db.Exec(query, post.Creater, post.Title, post.Description)
	if err != nil {
		return 0, fmt.Errorf("storage: create post: %w", err)
	}
	postId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("storage: create post: %w", err)
	}
	query = `UPDATE user SET posts = posts + 1 WHERE username = $1;`
	_, err = s.db.Exec(query, post.Creater)
	if err != nil {
		return 0, fmt.Errorf("storage: create post: %w", err)
	}
	query = `INSERT INTO post_category (postId, category) VALUES ($1, $2);`
	for _, oneCategory := range post.Category {
		_, err := s.db.Exec(query, postId, oneCategory)
		if err != nil {
			return 0, fmt.Errorf("storage: create post: %w", err)
		}
	}
	return int(postId), nil
}

func (s *PostStorage) GetAllPost() ([]models.Post, error) {
	queryPost := `SELECT * FROM post;`
	rows, err := s.db.Query(queryPost)
	if err != nil {
		return nil, fmt.Errorf("storage: get all post: %w", err)
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, fmt.Errorf("storage: get all post: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostStorage) GetPostById(postId int) (models.Post, error) {
	var post models.Post
	query := `SELECT * FROM post WHERE id = $1;`
	err := s.db.QueryRow(query, postId).Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments)
	if err != nil {
		return models.Post{}, fmt.Errorf("storage: get post by id: %w", err)
	}
	return post, nil
}

func (s *PostStorage) GetPostsByCategory(category string) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT * FROM post WHERE id IN (SELECT postId FROM post_category WHERE category=$1);`
	rows, err := s.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("storage: get post by category: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, fmt.Errorf("storage: get post by category: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostStorage) GetPostByTimeNew() ([]models.Post, error) {
	var posts []models.Post
	query := `select * from post ORDER by created_at DESC;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage: get post by time new: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, fmt.Errorf("storage: get post by time new: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostStorage) GetPostByTimeOld() ([]models.Post, error) {
	var posts []models.Post
	query := `select * from post ORDER by created_at;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage: get post by time old: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, fmt.Errorf("storage: get post by time old: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostStorage) GetPostByLikeMost() ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT * FROM post ORDER BY likes DESC;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage: get post by like most: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, fmt.Errorf("storage: get post by like most: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostStorage) GetPostByLikeLeast() ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT * FROM post ORDER BY likes;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("storage: get post by like least: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments); err != nil {
			return nil, fmt.Errorf("storage: get post by like least: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostStorage) GetSimilarPosts(postId int) ([]models.Post, error) {
	var posts []models.Post
	query := `select id, creater, title, description from post 
	WHERE id IN (SELECT postId FROM post_category 
	WHERE NOT postId = $1 AND category IN (SELECT category FROM post_category WHERE postid=$2));`
	rows, err := s.db.Query(query, postId, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get similar posts: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.Creater, &post.Title, &post.Description); err != nil {
			return nil, fmt.Errorf("storage: get similar posts: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostStorage) GetAllCategoryByPostId(postId int) ([]string, error) {
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

func (s *PostStorage) DeletePost(post models.Post) error {
	query := `UPDATE user 
	SET posts = posts - 1,
	likes = likes - $1,
	dislikes = dislikes - $2,
	comments = comments - $3
	WHERE username = $4;`
	_, err := s.db.Exec(query, post.Likes, post.Dislikes, post.Comments, post.Creater)
	if err != nil {
		return fmt.Errorf("storage: delete post: %w", err)
	}
	query = `DELETE FROM post WHERE id = $1;`
	_, err = s.db.Exec(query, post.Id)
	if err != nil {
		return fmt.Errorf("storage: delete post: %w", err)
	}
	return nil
}

func (s *PostStorage) ChangePost(newPost models.Post, postId int) error {
	query := `UPDATE post SET title = $1, description = $2 WHERE id = $3;`
	_, err := s.db.Exec(query, newPost.Title, newPost.Description, postId)
	if err != nil {
		return fmt.Errorf("storage: change post: %w", err)
	}
	query = `DELETE FROM post_category WHERE postId = $1;`
	_, err = s.db.Exec(query, postId)
	if err != nil {
		return fmt.Errorf("storage: change post: %w", err)
	}
	query = `INSERT INTO post_category (postId, category) VALUES ($1, $2);`
	for _, oneCategory := range newPost.Category {
		_, err = s.db.Exec(query, postId, oneCategory)
		if err != nil {
			return fmt.Errorf("storage: create post: %w", err)
		}
	}
	return nil
}

func (s *PostStorage) GetAllImagesByPostId(postId int) ([]string, error) {
	var images []string
	query := `SELECT image FROM post_images WHERE postId = $1;`
	rows, err := s.db.Query(query, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get all image by post id: %w", err)
	}
	for rows.Next() {
		var image string
		if err := rows.Scan(&image); err != nil {
			return nil, fmt.Errorf("storage: get all image by post id: %w", err)
		}
		images = append(images, image)
	}
	return images, nil
}

func (s *PostStorage) SaveImageForPost(postId int, imgPath string) error {
	query := `INSERT INTO post_images(postId, image) values ($1, $2);`
	_, err := s.db.Exec(query, postId, imgPath)
	if err != nil {
		return fmt.Errorf("storage: save image for post: %w", err)
	}
	return nil
}
