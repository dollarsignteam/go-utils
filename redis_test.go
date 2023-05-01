package utils_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
	"github.com/stretchr/testify/assert"

	"github.com/dollarsignteam/go-utils"
)

func TestRedisNew(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}
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
				URL: fmt.Sprintf("redis://%s", s.Addr()),
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
			URL: fmt.Sprintf("redis://%s", s.Addr()),
		})
		assert.EqualError(t, err, "mock error")
	})
}
