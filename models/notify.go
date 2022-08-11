package models

type Notify struct {
	Id          int
	From        string
	To          string
	Description string
	PostId      int
	CommentId   int
}
