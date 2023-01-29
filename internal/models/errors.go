package models

import "errors"

// user errors
var (
	ErrInvalidUser  = errors.New("invalid user")
	ErrUserNotFound = errors.New("user not found")
)

// reaction errors
var (
	ErrInvalidReaction = errors.New("invalid reaction")
)

// post errors
var (
	ErrInvalidFilter = errors.New("invalid filter")
)
