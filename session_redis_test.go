package utils_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
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
	_, err := h.Get("sid", 1, "gid")
	assert.ErrorIs(t, err, utils.ErrSessionNotFound)
}
