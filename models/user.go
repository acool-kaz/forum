package models

import "time"

type User struct {
	ID              int
	Email           string
	Username        string
	Password        string
	VerifyPassword  string
	CountOfPosts    int
	CountOfLikes    int
	CountOfComments int
	Token           string
	ExpiresAt       time.Time
}
