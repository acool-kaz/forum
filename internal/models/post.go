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
	Likes       uint          `json:"likes"`
	Dislikes    uint          `json:"dislikes"`
	Username    string        `json:"username"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Tags        []string      `json:"tags"`
	Images      []string      `json:"images"`
	Comments    []FullComment `json:"comments"`
	CreatedAt   time.Time     `json:"created_at"`
}

type postFilterCtx string

var (
	Tags   postFilterCtx = "tags"
	Filter postFilterCtx = "filter"
)
