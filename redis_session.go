package utils

import (
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Default redis session constants
const (
	DefaultRedisSessionKey = "session"
	DefaultRedisUserKey    = "user"
)

// RedisSessionConfig represents the configuration used to setup a Redis session manager.
type RedisSessionConfig struct {
	SessionKey             string        // The key used to store session data in Redis
	UserKey                string        // The key used to identify the user associated with a given session
	MultipleSessionPerUser bool          // Indicates whether multiple sessions are allowed per user
	Client                 *redis.Client // The Redis client used to communicate with the Redis server
}

// RedisSessionHandler represents a handler for Redis sessions.
type RedisSessionHandler struct {
	prefixKey              string
	multipleSessionPerUser bool
	client                 *redis.Client
}

// NewSessionHandler returns a new Redis session handler instance based on the provided configuration.
func (RedisUtil) NewSessionHandler(config RedisSessionConfig) *RedisSessionHandler {
	config.SessionKey = strings.TrimSpace(config.SessionKey)
	config.UserKey = strings.TrimSpace(config.UserKey)
	if config.SessionKey == "" {
		config.SessionKey = DefaultRedisSessionKey
	}
	if config.UserKey == "" {
		config.UserKey = DefaultRedisUserKey
	}
	return &RedisSessionHandler{
		prefixKey:              fmt.Sprintf("%s:%s", config.SessionKey, config.UserKey),
		multipleSessionPerUser: config.MultipleSessionPerUser,
		client:                 config.Client,
	}
}
