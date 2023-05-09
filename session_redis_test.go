package utils_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func createMockRedisClient(t *testing.T) (*miniredis.Miniredis, *utils.RedisClient) {
	s, url := createMockRedisServer(t)
	client, err := utils.Redis.New(utils.RedisConfig{
		URL: url,
	})
	if err != nil {
		t.Fatal(err)
	}
	return s, client
}

func TestRedis_NewSessionHandler(t *testing.T) {
	h := utils.Redis.NewSessionHandler(utils.SessionRedisConfig{})
	assert.IsType(t, &utils.SessionRedisHandler{}, h)
}

func TestSessionRedisHandler_Set(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := utils.Redis.NewSessionHandler(utils.SessionRedisConfig{
		SessionKey:             "s",
		UserKey:                "u",
		MultipleSessionPerUser: true,
		Client:                 client,
	})
	session := utils.Session{
		ID:      "sid",
		UserID:  1,
		GroupID: "gid",
		Data:    "foo",
	}

	t.Run("set", func(t *testing.T) {
		err := h.Set(session, time.Now().Add(1*time.Second).Unix())
		assert.Nil(t, err)
	})

	t.Run("get", func(t *testing.T) {
		expected := utils.Session{}
		err := client.GetStruct("s:u:gid:1:sid", &expected)
		assert.Nil(t, err)
		assert.Equal(t, session, expected)
	})
}

func TestSessionRedisHandler_Get(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := utils.Redis.NewSessionHandler(utils.SessionRedisConfig{
		SessionKey:             "s",
		UserKey:                "u",
		MultipleSessionPerUser: true,
		Client:                 client,
	})

	session := utils.Session{
		ID:      "foo",
		UserID:  1,
		GroupID: "bar",
		Data:    nil,
	}

	t.Run("success", func(t *testing.T) {
		err := h.Set(session, time.Now().Add(1*time.Second).Unix())
		assert.Nil(t, err)
		result, _ := h.Get("foo", 1, "bar")
		assert.Equal(t, session, result)
	})

	t.Run("invalid session", func(t *testing.T) {
		_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
		_ = client.SetStruct("s:u:bar:1:foo", utils.Session{}, 1*time.Second)
		_, err := h.Get("foo", 1, "bar")
		assert.ErrorIs(t, err, utils.ErrSessionInvalid)
	})

	t.Run("session not found", func(t *testing.T) {
		_, err := h.Get("sid", 1, "gid")
		assert.ErrorIs(t, err, utils.ErrSessionNotFound)
	})

	t.Run("redis error", func(t *testing.T) {
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			p.WriteError("mock error")
			return true
		})
		_, err := h.Get("sid", 1, "gid")
		assert.EqualError(t, err, "mock error")
	})

}
