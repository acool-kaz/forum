package models

type Info struct {
	User             User
	ProfileUser      User
	Posts            []FullPost
	SimilarPosts     []Post
	Post             FullPost
	Notifications    []Notification
	PostLikes        []string
	PostDislikes     []string
	Comments         []Comment
	CommentsLikes    map[int][]string
	CommentsDislikes map[int][]string
}
