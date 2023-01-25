package models

type Reaction struct {
	Id        uint `json:"id"`
	PostId    uint `json:"post_id"`
	CommentId uint `json:"comment_id"`
	UserId    uint `json:"user_id"`
	React     int  `json:"react"`
}

const (
	LikeReaction    = 1
	DislikeReaction = -1
)
