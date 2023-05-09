package utils

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
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

	t.Run("success", func(t *testing.T) {
		_ = s.Set("foo:bar", "baz")
		err := h.deleteSessionKeys("foo:*")
		assert.Nil(t, err)
		found := s.Exists("foo:var")
		assert.False(t, found)
	})

	t.Run("error", func(t *testing.T) {
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			p.WriteError("mock error")
			return true
		})
		err := h.deleteSessionKeys("foo:*")
		assert.EqualError(t, err, "mock error")
	})

	t.Run("pipe error", func(t *testing.T) {
		t.Skip()
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			if s1 == "SCAN" {
				return true
			}
			p.WriteError("mock error")
			return true
		})
		err := h.deleteSessionKeys("foo:*")
		assert.EqualError(t, err, "mock error")
	})
}
