package redis

import "time"

type Config struct {
	Addr              string
	ConnectionTimeout time.Duration
	QueryTimeout      time.Duration
	PostTTL           time.Duration
	SchemaTTL         time.Duration
}


