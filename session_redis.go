package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Default redis session constants
const (
	DefaultRedisSessionKey = "session"
	DefaultRedisUserKey    = "user"
)

const (
	sessionRedisLimitGetTx    = 100
	sessionRedisLimitDeleteTx = 1000
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionInvalid  = errors.New("invalid session")
)

// SessionRedisConfig represents the configuration used to setup a Redis session manager.
type SessionRedisConfig struct {
	SessionKey             string       // The key used to store session data in Redis
	UserKey                string       // The key used to identify the user associated with a given session
	MultipleSessionPerUser bool         // Indicates whether multiple sessions are allowed per user
	Client                 *RedisClient // The Redis client used to communicate with the Redis server
}

// SessionRedisGetKeyParam represents parameters used to generate Redis keys for session storage.
type SessionRedisGetKeyParam struct {
	SessionID string
	UserID    string
	GroupID   string
}

// Session represents a single session stored in Redis.
type Session struct {
	ID      string `json:"sid"`
	UserID  int64  `json:"uid"`
	GroupID string `json:"gid,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type SessionHandler interface {
	Set(session Session, expiresAt int64) error
	Get(sessionId string, userId int64, groupId string) (Session, error)
	ListByUserID(userId int64) ([]Session, error)
	ListByGroupID(groupId string) ([]Session, error)
	Exists(sessionId string, userId int64, groupId string) (bool, error)
	Count(uniqueByUser bool) (int, error)
	CountByUserID(userId int64, uniqueByUser bool) (int, error)
	CountByGroupID(groupId string, uniqueByUser bool) (int, error)
	Delete(sessionId string, userId int64, groupId string) error
	DeleteByUserID(userId int64) error
	DeleteByGroupID(groupId string) error
	DeleteAll() error
}

// SessionRedisHandler represents a handler for Redis sessions.
type SessionRedisHandler struct {
	multipleSessionPerUser bool
	prefixKey              string
	client                 *RedisClient
}

// NewSessionHandler returns a new Redis session handler
// instance based on the provided configuration.
func (RedisUtil) NewSessionHandler(config SessionRedisConfig) SessionHandler {
	config.SessionKey = strings.TrimSpace(config.SessionKey)
	config.UserKey = strings.TrimSpace(config.UserKey)
	if config.SessionKey == "" {
		config.SessionKey = DefaultRedisSessionKey
	}
	if config.UserKey == "" {
		config.UserKey = DefaultRedisUserKey
	}
	return &SessionRedisHandler{
		multipleSessionPerUser: config.MultipleSessionPerUser,
		prefixKey:              fmt.Sprintf("%s:%s", config.SessionKey, config.UserKey),
		client:                 config.Client,
	}
}

// Set stores the provided session in Redis.
func (h *SessionRedisHandler) Set(s Session, expiresAt int64) error {
	key := h.getKey(SessionRedisGetKeyParam{
		SessionID: s.ID,
		UserID:    strconv.FormatInt(s.UserID, 10),
		GroupID:   s.GroupID,
	})
	ttl := time.Duration(Max(1, expiresAt-time.Now().Unix())) * time.Second
	return h.client.SetStruct(key, s, time.Duration(ttl)*time.Second)
}

// Get retrieves the session with the specified ID for the given user and group.
func (h *SessionRedisHandler) Get(sessionId string, userId int64, groupId string) (Session, error) {
	key := h.getKey(SessionRedisGetKeyParam{
		SessionID: sessionId,
		UserID:    strconv.FormatInt(userId, 10),
		GroupID:   groupId,
	})
	session := Session{}
	if err := h.client.GetStruct(key, &session); err != nil {
		if err == redis.Nil {
			return session, ErrSessionNotFound
		}
		return session, err
	}
	if session.UserID != userId || session.ID != sessionId {
		return session, ErrSessionInvalid
	}
	return session, nil
}

// ListByUserID retrieves all sessions for the given user.
func (h *SessionRedisHandler) ListByUserID(userId int64) ([]Session, error) {
	key := h.getKey(SessionRedisGetKeyParam{
		SessionID: "*",
		UserID:    strconv.FormatInt(userId, 10),
		GroupID:   "*",
	})
	return h.find(key, func(s *Session) bool {
		return s.UserID == userId
	})
}

// ListByGroupID retrieves all sessions for the given group.
func (h *SessionRedisHandler) ListByGroupID(groupId string) ([]Session, error) {
	if groupId == "*" {
		return []Session{}, nil
	}
	key := h.getKey(SessionRedisGetKeyParam{
		SessionID: "*",
		UserID:    "*",
		GroupID:   groupId,
	})
	return h.find(key, func(s *Session) bool {
		return s.GroupID == groupId
	})
}

// Exists checks if the session with the specified ID for the given user and group exists.
func (h *SessionRedisHandler) Exists(sessionId string, userId int64, groupId string) (bool, error) {
	key := h.getKey(SessionRedisGetKeyParam{
		SessionID: sessionId,
		UserID:    strconv.FormatInt(userId, 10),
		GroupID:   groupId,
	})
	res, err := h.client.Exists(context.TODO(), key).Result()
	return res != 0, err
}

// Count returns the number of stored sessions.
func (h *SessionRedisHandler) Count(uniqueByUser bool) (int, error) {
	key := fmt.Sprintf("%s:*", h.prefixKey)
	return h.countByKey(key, uniqueByUser)
}

// CountByUserID returns the number of sessions associated with the given user.
func (h *SessionRedisHandler) CountByUserID(userId int64, uniqueByUser bool) (int, error) {
	key := fmt.Sprintf("%s:*:%d:*", h.prefixKey, userId)
	return h.countByKey(key, uniqueByUser)
}

// CountByGroupID returns the number of sessions associated with the given group.
func (h *SessionRedisHandler) CountByGroupID(groupId string, uniqueByUser bool) (int, error) {
	key := fmt.Sprintf("%s:%s:*", h.prefixKey, groupId)
	return h.countByKey(key, uniqueByUser)
}

// Delete removes the session with the specified ID for the given user and group from Redis.
func (h *SessionRedisHandler) Delete(sessionId string, userId int64, groupId string) error {
	key := h.getKey(SessionRedisGetKeyParam{
		SessionID: sessionId,
		UserID:    strconv.FormatInt(userId, 10),
		GroupID:   groupId,
	})
	return h.client.Del(context.TODO(), key).Err()
}

// DeleteByUserID removes all sessions associated with the given user from Redis.
func (h *SessionRedisHandler) DeleteByUserID(userId int64) error {
	key := fmt.Sprintf("%s:*:%d:*", h.prefixKey, userId)
	return h.deleteSessionKeys(key)
}

// DeleteByGroupID removes all sessions associated with the given group from Redis.
func (h *SessionRedisHandler) DeleteByGroupID(groupId string) error {
	key := fmt.Sprintf("%s:%s:*", h.prefixKey, groupId)
	return h.deleteSessionKeys(key)
}

// DeleteAll removes all sessions from Redis.
func (h *SessionRedisHandler) DeleteAll() error {
	key := fmt.Sprintf("%s:*", h.prefixKey)
	return h.deleteSessionKeys(key)
}

// getKey returns the Redis key for the given session.
func (h *SessionRedisHandler) getKey(param SessionRedisGetKeyParam) string {
	var builder strings.Builder
	builder.WriteString(h.prefixKey)
	builder.WriteByte(':')
	builder.WriteString(param.GroupID)
	builder.WriteByte(':')
	builder.WriteString(param.UserID)
	if h.multipleSessionPerUser {
		builder.WriteByte(':')
		builder.WriteString(param.SessionID)
	}
	return builder.String()
}

// find searches Redis for all keys matching a given pattern and returns a slice of all sessions
// that pass a verification function. Uses the SCAN command and Redis pipelines.
func (h *SessionRedisHandler) find(key string, verifyFunc func(*Session) bool) ([]Session, error) {
	var cursor uint64
	ctx := context.TODO()
	sessions := []Session{}
	for {
		var err error
		var keys []string
		keys, cursor, err = h.client.Scan(ctx, cursor, key, sessionRedisLimitGetTx).Result()
		if err != nil {
			return nil, err
		}
		pipe := h.client.Pipeline()
		for _, k := range keys {
			pipe.Get(ctx, k)
		}
		cmds, err := pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
		for _, cmd := range cmds {
			res, err := cmd.(*redis.StringCmd).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				return nil, err
			}
			var session Session
			err = json.Unmarshal([]byte(res), &session)
			if err != nil || verifyFunc(&session) {
				continue
			}
			sessions = append(sessions, session)
		}
		if cursor == 0 {
			break
		}
	}
	return sessions, nil
}

// countByKey counts the number of session keys based on the given 'key'.
// If 'uniqueByUser' is set to true, it only counts unique sessions by user.
// It returns the total count or the unique count and any potential errors.
func (h *SessionRedisHandler) countByKey(key string, uniqueByUser bool) (int, error) {
	var count int
	if uniqueByUser {
		m := make(map[string]struct{})
		err := h.scanSessionKeys(key, func(k string) {
			s := strings.Split(k, ":")
			if len(s) > 3 {
				m[s[3]] = struct{}{}
			}
		})
		return len(m), err
	}
	err := h.scanSessionKeys(key, func(k string) {
		count++
	})
	return count, err
}

// scans Redis session keys using SCAN command and processes them one by one.
func (h *SessionRedisHandler) scanSessionKeys(key string, processKeyFunc func(string)) error {
	var cursor uint64
	ctx := context.TODO()
	for {
		var err error
		var keys []string
		keys, cursor, err = h.client.Scan(ctx, cursor, key, sessionRedisLimitGetTx).Result()
		if err != nil {
			return err
		}
		for _, k := range keys {
			processKeyFunc(k)
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

// deleteSessionKeys scans and deletes Redis session keys
// using the SCAN command and Redis transactions.
func (h *SessionRedisHandler) deleteSessionKeys(key string) error {
	ctx := context.TODO()
	for {
		keys, cursor, err := h.client.Scan(ctx, 0, key, sessionRedisLimitDeleteTx).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			pipe := h.client.TxPipeline()
			for _, k := range keys {
				pipe.Del(ctx, k)
			}
			if _, err = pipe.Exec(ctx); err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}
