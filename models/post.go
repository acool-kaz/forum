package models

import "time"

type Post struct {
	Id          uint      `json:"id"`
	UserId      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Tags        string    `json:"tags"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type FullPost struct {
	Id          uint          `json:"id"`
	Username    string        `json:"username"`
	Title       string        `json:"title"`
	Tags        []string      `json:"tags"`
	Description string        `json:"description"`
	Likes       uint          `json:"likes"`
	Dislikes    uint          `json:"dislikes"`
	Comments    []FullComment `json:"comments"`
	Images      []string      `json:"images"`
	CreatedAt   time.Time     `json:"created_at"`
}
