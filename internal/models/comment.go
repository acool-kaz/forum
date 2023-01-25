package models

type Comment struct {
	Id     uint   `json:"id"`
	PostId uint   `json:"post_id"`
	UserId uint   `json:"user_id"`
	Text   string `json:"text"`
}

type FullComment struct {
	Id       uint   `json:"id"`
	Likes    uint   `json:"likes"`
	Dislikes uint   `json:"dislikes"`
	Username string `json:"username"`
	Text     string `json:"text"`
}