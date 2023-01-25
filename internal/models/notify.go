package models

type Notification struct {
	Id          int
	From        string
	To          string
	Description string
	PostId      int
	CommentId   int
}
