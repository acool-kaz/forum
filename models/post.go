package models

import "time"

type Post struct {
	Id          int
	Creater     string
	Category    []string
	Title       string
	Description string
	CreatedAt   time.Time
	Likes       int
	Dislikes    int
	Comments    int
}
