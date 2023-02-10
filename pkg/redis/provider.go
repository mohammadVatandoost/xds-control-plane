package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func ProvideRedisClient(addr string, dialTimeout time.Duration,
	hooks []redis.Hook) *Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr:        addr,
		DialTimeout: dialTimeout,
	})

	for _, hook := range hooks {
		redisClient.AddHook(hook)
	}

	client := NewClient(redisClient)

	return client
}

func ProvideRedisSentinelClient(addr string, masterName string, dialTimeout time.Duration,
	hooks []redis.Hook) (*Client, error) {
	sentinel := redis.NewSentinelClient(&redis.Options{
		Addr: addr,
	})

	address, err := sentinel.GetMasterAddrByName(context.Background(), masterName).Result()

	if err != nil {
		return nil, err
	}

	return ProvideRedisClient(address[0], dialTimeout, hooks), nil
}
