package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis utility instance
var Redis redisUtil

// redisUtil is a utility struct for working with Redis
type redisUtil struct{}

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
func (redisUtil) New(config RedisConfig) (*RedisClient, error) {
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

// SetStruct sets the value for a key in Redis
// The value is marshalled to JSON, and an optional expiration time can be set
func (r *RedisClient) SetStruct(key string, value any, expiration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(context.TODO(), key, string(b), expiration).Err()
}

// SetNXStruct sets the value for a key in Redis if the key does not already exist
// The value is marshalled to JSON, and an optional expiration time can be set
// Returns a boolean indicating whether the key was set, and an error (if any)
func (r *RedisClient) SetNXStruct(key string, value any, expiration time.Duration) (bool, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.Client.SetNX(context.TODO(), key, string(b), expiration).Result()
}

// GetStruct retrieves the value of a key as a JSON-encoded struct
// and unmarshal it into a provided result variable
func (r *RedisClient) GetStruct(key string, result any) error {
	val, err := r.Client.Get(context.TODO(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), result)
}

// JSONSet is a convenience wrapper around SetStruct
func (r *RedisClient) JSONSet(key string, value any, expiration time.Duration) error {
	return r.SetStruct(key, value, expiration)
}

// JSONGet is a convenience wrapper around GetStruct
func (r *RedisClient) JSONGet(key string, result any) error {
	return r.GetStruct(key, result)
}
