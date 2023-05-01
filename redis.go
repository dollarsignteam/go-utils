package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis utility instance
var Redis RedisUtil

// RedisUtil is a utility struct for working with Redis
type RedisUtil struct{}

// redisPingTTL is the timeout duration for a Ping request to Redis
var redisPingTTL = 5 * time.Second

// RedisConfig specifies configuration options for connecting to a Redis server
type RedisConfig struct {
	URL string
}

// RedisClient is a wrapper around go-redis Client type
type RedisClient struct {
	*redis.Client
}

// New creates a new Redis client based on the provided RedisConfig
func (RedisUtil) New(config RedisConfig) (*RedisClient, error) {
	opt, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), redisPingTTL)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{client}, err
}
