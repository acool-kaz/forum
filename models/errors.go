package models

import "errors"

var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUserNotFound = errors.New("user not found")
)
