package entity

import "errors"

// Domain errors
var (
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidChatID      = errors.New("invalid chat ID")
	ErrEmptyContent       = errors.New("memory content cannot be empty")
	ErrMemoryNotFound     = errors.New("memory not found")
	ErrUnauthorized       = errors.New("unauthorized access to memory")
	ErrInvalidSearchQuery = errors.New("invalid search query")
)
