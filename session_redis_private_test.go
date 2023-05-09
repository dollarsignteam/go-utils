package utils

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func createMockRedisClient(t *testing.T) (*miniredis.Miniredis, *RedisClient) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	url := fmt.Sprintf("redis://%s", s.Addr())
	client, err := Redis.New(RedisConfig{
		URL: url,
	})
	if err != nil {
		t.Fatal(err)
	}
	return s, client
}

func createHandler(client *RedisClient) *SessionRedisHandler {
	return &SessionRedisHandler{
		multipleSessionPerUser: true,
		prefixKey:              "s:u",
		client:                 client,
	}
}

func TestSessionRedisHandler_deleteSessionKeys(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := createHandler(client)
	_ = s.Set("foo:bar", "baz")
	err := h.deleteSessionKeys("foo:*")
	assert.Nil(t, err)
	found := s.Exists("foo:var")
	assert.False(t, found)
}
