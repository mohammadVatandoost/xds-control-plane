package redis

import "errors"

var (
	ErrNotFound     = errors.New("key not found")
	ErrNotAvailable = errors.New("not available")
)

