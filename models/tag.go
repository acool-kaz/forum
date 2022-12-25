package models

type Tag struct {
	Id     uint   `json:"id"`
	PostId uint   `json:"post_id"`
	Name   string `json:"name"`
}
