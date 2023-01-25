package models

import "time"

type Session struct {
	Id        uint      `json:"id"`
	UserId    uint      `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
