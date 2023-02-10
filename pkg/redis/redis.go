package redis

import (
	"context"
	"time"
)

type Cmdable interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
}

