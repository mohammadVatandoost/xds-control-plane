package redis

import "time"

type Config struct {
	Address           string
	ConnectionTimeout time.Duration
}
