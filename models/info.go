package models

type Info struct {
	Posts            []Post
	SimilarPosts     []Post
	Post             Post
	PostLikes        []string
	PostDislikes     []string
	User             User
	Comments         []Comment
	CommentsLikes    map[int][]string
	CommentsDislikes map[int][]string
}
