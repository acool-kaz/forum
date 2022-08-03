package models

type Comment struct {
	Id       int
	PostId   int
	Creater  string
	Text     string
	Likes    int
	Dislikes int
}
