package models

type Info struct {
	User             User
	ProfileUser      User
	Posts            []Post
	SimilarPosts     []Post
	Post             Post
	Notifications    []Notification
	PostLikes        []string
	PostDislikes     []string
	Comments         []Comment
	CommentsLikes    map[int][]string
	CommentsDislikes map[int][]string
}
