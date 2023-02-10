package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	base redis.Cmdable
}

func NewClient(base redis.Cmdable) *Client {
	return &Client{
		base: base,
	}
}

func (c *Client) get(ctx context.Context, key string) (string, error) {

	res, err := c.base.Get(ctx, key).Result()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	result, err := c.get(ctx, key)

	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "", ErrNotAvailable
	}

	return result, errors.WithStack(err)
}

func (c *Client) set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {

	res, err := c.base.Set(ctx, key, value, expiration).Result()

	// grpc ok errors
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	result, err := c.set(ctx, key, value, expiration)

	if errors.Is(err, context.DeadlineExceeded) {
		return "", ErrNotAvailable
	}

	return result, errors.WithStack(err)
}
