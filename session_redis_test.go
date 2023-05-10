package utils

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func createMockRedisClient(t *testing.T) (*miniredis.Miniredis, *RedisClient) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	client, err := Redis.New(RedisConfig{
		URL: fmt.Sprintf("redis://%s", s.Addr()),
	})
	if err != nil {
		t.Fatal(err)
	}
	return s, client
}

func TestSession_Scan(t *testing.T) {
	tests := []struct {
		name        string
		session     Session
		dest        any
		expectedErr error
	}{
		{
			name: "valid data",
			session: Session{
				Meta: SessionMeta{},
				Data: map[string]interface{}{
					"name": "John",
					"age":  30,
				},
			},
			dest: &struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{},
			expectedErr: nil,
		},
		{
			name: "invalid data",
			session: Session{
				Meta: SessionMeta{},
				Data: "invalid json",
			},
			dest:        new(int),
			expectedErr: errors.New("json: cannot unmarshal string into Go value of type int"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.session.Scan(test.dest)
			if test.expectedErr != nil {
				assert.EqualError(t, err, test.expectedErr.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSessionMeta_param(t *testing.T) {
	tests := []struct {
		name     string
		meta     SessionMeta
		expected sessionKeyParam
	}{
		{
			name: "no group ID",
			meta: SessionMeta{
				ID:     "1234",
				UserID: 5678,
			},
			expected: sessionKeyParam{
				ID:      "1234",
				userID:  "5678",
				groupID: "",
			},
		},
		{
			name: "with group ID",
			meta: SessionMeta{
				ID:      "abcd",
				UserID:  9876,
				GroupID: "xyz",
			},
			expected: sessionKeyParam{
				ID:      "abcd",
				userID:  "9876",
				groupID: "xyz",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.meta.param()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestRedis_NewSessionHandler(t *testing.T) {
	tests := []struct {
		name     string
		config   SessionRedisConfig
		expected SessionHandler
	}{
		{
			name: "empty config",
			config: SessionRedisConfig{
				Client: &RedisClient{},
			},
			expected: &SessionRedisHandler{
				multipleSessionPerUser: false,
				prefixKey:              fmt.Sprintf("%s:%s", DefaultRedisSessionKey, DefaultRedisUserKey),
				client:                 &RedisClient{},
			},
		},
		{
			name: "session key only",
			config: SessionRedisConfig{
				SessionKey: "sess",
				Client:     &RedisClient{},
			},
			expected: &SessionRedisHandler{
				multipleSessionPerUser: false,
				prefixKey:              fmt.Sprintf("sess:%s", DefaultRedisUserKey),
				client:                 &RedisClient{},
			},
		},
		{
			name: "user key only",
			config: SessionRedisConfig{
				UserKey: "use",
				Client:  &RedisClient{},
			},
			expected: &SessionRedisHandler{
				multipleSessionPerUser: false,
				prefixKey:              fmt.Sprintf("%s:use", DefaultRedisSessionKey),
				client:                 &RedisClient{},
			},
		},
		{
			name: "multiple sessions per user",
			config: SessionRedisConfig{
				Client:                 &RedisClient{},
				MultipleSessionPerUser: true,
			},
			expected: &SessionRedisHandler{
				multipleSessionPerUser: true,
				prefixKey:              fmt.Sprintf("%s:%s", DefaultRedisSessionKey, DefaultRedisUserKey),
				client:                 &RedisClient{},
			},
		},
		{
			name: "custom keys",
			config: SessionRedisConfig{
				SessionKey:             "sess",
				UserKey:                "use",
				Client:                 &RedisClient{},
				MultipleSessionPerUser: true,
			},
			expected: &SessionRedisHandler{
				multipleSessionPerUser: true,
				prefixKey:              "sess:use",
				client:                 &RedisClient{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Redis.NewSessionHandler(test.config)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestSessionRedisHandler_Set(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: true,
		Client:                 client,
	})
	session := Session{
		Meta: SessionMeta{
			ID:      "foo",
			UserID:  1,
			GroupID: "bar",
		},
		Data: "baz",
	}

	t.Run("set", func(t *testing.T) {
		err := h.Set(session, time.Now().Add(1*time.Second).Unix())
		assert.Nil(t, err)
	})

	t.Run("get", func(t *testing.T) {
		expected := Session{}
		err := client.GetStruct("session:user:bar:1:foo", &expected)
		assert.Nil(t, err)
		assert.Equal(t, session, expected)
	})
}

func BenchmarkSession_Scan(b *testing.B) {
	session := Session{
		Meta: SessionMeta{},
		Data: map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}
	dest := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}
	for i := 0; i < b.N; i++ {
		err := session.Scan(dest)
		assert.Nil(b, err)
	}
}
