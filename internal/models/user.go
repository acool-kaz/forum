package models

type User struct {
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	VerifyPassword string `json:"verify_password"`
}

type userCtx string

var (
	UserId   userCtx = "user_id"
	Username userCtx = "username"
	Email    userCtx = "email"
)
