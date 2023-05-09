package utils_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func createMockRedisServer(t *testing.T) (*miniredis.Miniredis, string) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	return s, fmt.Sprintf("redis://%s", s.Addr())
}

func TestRedisNew(t *testing.T) {
	s, url := createMockRedisServer(t)
	defer s.Close()
	tests := []struct {
		name          string
		config        utils.RedisConfig
		expectedError error
	}{
		{
			name: "invalid url",
			config: utils.RedisConfig{
				URL: "invalid://url",
			},
			expectedError: errors.New("redis: invalid URL scheme: invalid"),
		},
		{
			name: "success",
			config: utils.RedisConfig{
				URL: url,
			},
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := utils.Redis.New(test.config)
			if test.expectedError != nil {
				assert.Nil(t, result)
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				assert.NotNil(t, result)
			}
		})
	}

	t.Run("ping failed", func(t *testing.T) {
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			p.WriteError("mock error")
			return true
		})
		_, err := utils.Redis.New(utils.RedisConfig{
			URL: url,
		})
		assert.EqualError(t, err, "mock error")
	})
}

type testRedisStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSetAndGetStruct(t *testing.T) {
	s, url := createMockRedisServer(t)
	defer s.Close()
	client, err := utils.Redis.New(utils.RedisConfig{
		URL: url,
	})
	if err != nil {
		t.Fatalf("error creating Redis client: %v", err)
	}

	t.Run("marshal failed", func(t *testing.T) {
		err := client.SetStruct("foo", make(chan int), 0)
		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		s.FlushDB()
		key := "foo"
		testData := testRedisStruct{
			Name: "Alice",
			Age:  30,
		}
		err = client.SetStruct("foo", testData, time.Minute)
		if ok := assert.Nil(t, err); ok {
			result := testRedisStruct{}
			err = client.GetStruct(key, &result)
			assert.Nil(t, err)
			assert.Equal(t, testData, result)
		}
	})

	t.Run("data not found", func(t *testing.T) {
		s.FlushDB()
		result := testRedisStruct{}
		err := client.GetStruct("foo", &result)
		assert.EqualError(t, err, "redis: nil")
	})
}

func TestSetNXStruct(t *testing.T) {
	s, url := createMockRedisServer(t)
	defer s.Close()
	client, err := utils.Redis.New(utils.RedisConfig{
		URL: url,
	})
	if err != nil {
		t.Fatalf("error creating Redis client: %v", err)
	}

	t.Run("marshal failed", func(t *testing.T) {
		_, err := client.SetNXStruct("foo", make(chan int), 0)
		assert.Error(t, err)
	})

	t.Run("exist data", func(t *testing.T) {
		s.FlushDB()
		key := "foo"
		ok, err := client.SetNXStruct(key, true, 0)
		assert.True(t, ok)
		assert.Nil(t, err)
		ok, err = client.SetNXStruct(key, nil, 0)
		assert.False(t, ok)
		assert.Nil(t, err)
	})
}

func TestJSONSetAndGet(t *testing.T) {
	s, url := createMockRedisServer(t)
	defer s.Close()
	client, err := utils.Redis.New(utils.RedisConfig{
		URL: url,
	})
	if err != nil {
		t.Fatalf("error creating Redis client: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		s.FlushDB()
		key := "foo"
		err := client.JSONSet(key, true, 0)
		assert.Nil(t, err)
		result := new(bool)
		err = client.JSONGet(key, result)
		assert.True(t, *result)
		assert.Nil(t, err)
	})
}
