package utils

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
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

func TestSessionRedisHandler_Get(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		SessionKey:             "s",
		UserKey:                "u",
		MultipleSessionPerUser: true,
		Client:                 client,
	})
	meta := SessionMeta{
		ID:      "foo",
		UserID:  1,
		GroupID: "bar",
	}
	session := Session{
		Meta: meta,
		Data: nil,
	}

	t.Run("success", func(t *testing.T) {
		err := h.Set(session, time.Now().Add(1*time.Second).Unix())
		assert.Nil(t, err)
		result, _ := h.Get(meta)
		assert.Equal(t, session, result)
	})

	t.Run("invalid session", func(t *testing.T) {
		_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
		_ = client.SetStruct("s:u:bar:1:foo", Session{}, 1*time.Second)
		_, err := h.Get(meta)
		assert.ErrorIs(t, err, ErrSessionInvalid)
	})

	t.Run("session not found", func(t *testing.T) {
		s.FlushDB()
		_, err := h.Get(meta)
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})

	t.Run("redis error", func(t *testing.T) {
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			p.WriteError("mock error")
			return true
		})
		_, err := h.Get(meta)
		assert.EqualError(t, err, "mock error")
	})
}

func TestSessionRedisHandler_ListByUserID(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: false,
		Client:                 client,
	})
	session := Session{
		Meta: SessionMeta{
			ID:      "foo",
			UserID:  1,
			GroupID: "bar",
		},
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
	expected := []Session{session}
	result, _ := h.ListByUserID(1)
	assert.Equal(t, expected, result)
}

func TestSessionRedisHandler_ListByGroupID(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: false,
		Client:                 client,
	})
	session := Session{
		Meta: SessionMeta{
			ID:      "foo",
			UserID:  1,
			GroupID: "bar",
		},
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())

	t.Run("group == *", func(t *testing.T) {
		expected := []Session{}
		result, _ := h.ListByGroupID("*")
		assert.Equal(t, expected, result)
	})

	t.Run("list group", func(t *testing.T) {
		expected := []Session{session}
		result, _ := h.ListByGroupID("bar")
		assert.Equal(t, expected, result)
	})
}

func TestSessionRedisHandler_Exists(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: false,
		Client:                 client,
	})
	meta := SessionMeta{
		ID:      "foo",
		UserID:  1,
		GroupID: "bar",
	}
	session := Session{
		Meta: meta,
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
	result, _ := h.Exists(meta)
	assert.True(t, result)
}

func TestSessionRedisHandler_Count(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: true,
		Client:                 client,
	})
	session1 := Session{
		Meta: SessionMeta{
			ID:      "foo1",
			UserID:  1,
			GroupID: "bar",
		},
		Data: nil,
	}
	session2 := Session{
		Meta: SessionMeta{
			ID:      "foo2",
			UserID:  1,
			GroupID: "bar",
		},
		Data: nil,
	}
	_ = h.Set(session1, time.Now().Add(1*time.Second).Unix())
	_ = h.Set(session2, time.Now().Add(1*time.Second).Unix())

	t.Run("unique by user", func(t *testing.T) {
		expected := 1
		result, _ := h.Count(true)
		assert.Equal(t, expected, result)
	})

	t.Run("not unique by user", func(t *testing.T) {
		expected := 2
		result, _ := h.Count(false)
		assert.Equal(t, expected, result)
	})
}

func TestSessionRedisHandler_CountByUserID(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: false,
		Client:                 client,
	})
	session1 := Session{
		Meta: SessionMeta{
			ID:      "foo1",
			UserID:  1,
			GroupID: "bar",
		},
		Data: nil,
	}
	session2 := Session{
		Meta: SessionMeta{
			ID:      "foo2",
			UserID:  1,
			GroupID: "bar",
		},
		Data: nil,
	}
	_ = h.Set(session1, time.Now().Add(1*time.Second).Unix())
	_ = h.Set(session2, time.Now().Add(1*time.Second).Unix())

	t.Run("unique by user", func(t *testing.T) {
		expected := 1
		result, _ := h.CountByUserID(1, true)
		assert.Equal(t, expected, result)
	})

	t.Run("not unique by user", func(t *testing.T) {
		expected := 1
		result, _ := h.CountByUserID(1, false)
		assert.Equal(t, expected, result)
	})
}

func TestSessionRedisHandler_CountByGroupID(t *testing.T) {
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
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())

	t.Run("unique by user", func(t *testing.T) {
		expected := 0
		result, _ := h.CountByGroupID("foo", true)
		assert.Equal(t, expected, result)
	})

	t.Run("not unique by user", func(t *testing.T) {
		expected := 1
		result, _ := h.CountByGroupID("bar", false)
		assert.Equal(t, expected, result)
	})
}

func TestSessionRedisHandler_Delete(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: false,
		Client:                 client,
	})
	meta := SessionMeta{
		ID:      "foo",
		UserID:  1,
		GroupID: "bar",
	}
	session := Session{
		Meta: meta,
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
	_ = h.Delete(meta)
	_, err := h.Get(meta)
	assert.ErrorIs(t, err, ErrSessionNotFound)
}

func TestSessionRedisHandler_DeleteByUserID(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: true,
		Client:                 client,
	})
	meta := SessionMeta{
		ID:      "foo",
		UserID:  1,
		GroupID: "bar",
	}
	session := Session{
		Meta: meta,
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
	_ = h.DeleteByUserID(1)
	_, err := h.Get(meta)
	assert.ErrorIs(t, err, ErrSessionNotFound)
}

func TestSessionRedisHandler_DeleteByGroupID(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: true,
		Client:                 client,
	})
	meta := SessionMeta{
		ID:      "foo",
		UserID:  1,
		GroupID: "bar",
	}
	session := Session{
		Meta: meta,
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
	_ = h.DeleteByGroupID("bar")
	_, err := h.Get(meta)
	assert.ErrorIs(t, err, ErrSessionNotFound)
}

func TestSessionRedisHandler_DeleteAll(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := Redis.NewSessionHandler(SessionRedisConfig{
		MultipleSessionPerUser: true,
		Client:                 client,
	})
	meta := SessionMeta{
		ID:      "foo",
		UserID:  1,
		GroupID: "bar",
	}
	session := Session{
		Meta: meta,
		Data: nil,
	}
	_ = h.Set(session, time.Now().Add(1*time.Second).Unix())
	_ = h.DeleteAll()
	count, _ := h.Count(false)
	assert.Equal(t, 0, count)
}

func TestSessionRedisHandler_find(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s, client := createMockRedisClient(t)
		defer s.Close()
		h := SessionRedisHandler{
			prefixKey: "session:user",
			client:    client,
		}
		session1 := Session{
			Meta: SessionMeta{
				ID:      "foo1",
				UserID:  1,
				GroupID: "bar",
			},
			Data: nil,
		}
		session2 := Session{
			Meta: SessionMeta{
				ID:      "foo2",
				UserID:  2,
				GroupID: "bar",
			},
			Data: nil,
		}
		_ = h.Set(session1, time.Now().Add(1*time.Second).Unix())
		_ = h.Set(session2, time.Now().Add(1*time.Second).Unix())
		result, _ := h.find("session:*", func(s *Session) bool {
			return s.Meta.UserID == 1
		})
		assert.Equal(t, session1, result[0])
	})

	t.Run("error", func(t *testing.T) {
		s, client := createMockRedisClient(t)
		defer s.Close()
		h := SessionRedisHandler{client: client}
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			p.WriteError("mock error")
			return true
		})
		_, err := h.find("foo:*", func(s *Session) bool {
			return true
		})
		assert.EqualError(t, err, "mock error")
	})

	t.Run("pipe error", func(t *testing.T) {
		s, client := createMockRedisClient(t)
		defer s.Close()
		h := SessionRedisHandler{
			prefixKey: "session:user",
			client:    client,
		}
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			if s1 == "SCAN" {
				p.WriteLen(2)
				p.WriteBulk("0")
				p.WriteLen(1)
				p.WriteBulk("foo:bar")
				return true
			}
			p.WriteError("pipe error")
			return true
		})
		_, err := h.find("foo:*", func(s *Session) bool {
			return true
		})
		assert.EqualError(t, err, "pipe error")
	})
}

func TestSessionRedisHandler_scanSessionKeys(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := SessionRedisHandler{client: client}

	t.Run("success", func(t *testing.T) {
		_ = s.Set("foo:bar", "baz")
		_ = h.scanSessionKeys("foo:*", func(s string) {
			assert.Equal(t, "foo:bar", s)
		})
	})

	t.Run("error", func(t *testing.T) {
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			p.WriteError("mock error")
			return true
		})
		err := h.scanSessionKeys("foo:*", func(s string) {})
		assert.EqualError(t, err, "mock error")
	})
}

func TestSessionRedisHandler_deleteSessionKeys(t *testing.T) {
	s, client := createMockRedisClient(t)
	defer s.Close()
	h := SessionRedisHandler{client: client}

	t.Run("success", func(t *testing.T) {
		_ = s.Set("foo:bar", "baz")
		_ = h.deleteSessionKeys("foo:*")
		result := s.Exists("foo:bar")
		assert.False(t, result)
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
		s, client := createMockRedisClient(t)
		defer s.Close()
		h := SessionRedisHandler{client: client}
		s.Server().SetPreHook(func(p *server.Peer, s1 string, s2 ...string) bool {
			if s1 == "SCAN" {
				p.WriteLen(2)
				p.WriteBulk("0")
				p.WriteLen(1)
				p.WriteBulk("foo:bar")
				return true
			}
			p.WriteError("pipe error")
			return true
		})
		err := h.deleteSessionKeys("foo:*")
		assert.EqualError(t, err, "pipe error")
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
