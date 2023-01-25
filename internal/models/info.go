package models

type Info struct {
	User          User
	ProfileUser   User
	Posts         []FullPost
	Post          FullPost
	SimilarPosts  []Post
	Notifications []Notification
}
