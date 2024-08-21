package utils

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// Default keys for storing session and user data in Redis.
const (
	DefaultRedisSessionKey = "session"
	DefaultRedisUserKey    = "user"
)

// Define limits for Redis transactions when retrieving or deleting session data.
const (
	sessionRedisLimitGetTx    = 100  // Max Redis keys fetched per tx
	sessionRedisLimitDeleteTx = 1000 // Max Redis keys deleted per tx
)

var (
	ErrSessionNotFound = errors.New("session not found") // Error for when session is not found.
	ErrSessionInvalid  = errors.New("invalid session")   // Error for when session is invalid.
)

// SessionRedisConfig is used to configure session data stored in Redis.
type SessionRedisConfig struct {
	SessionKey             string       // key for storing session data in Redis
	UserKey                string       // key for storing user-session mappings in Redis
	MultipleSessionPerUser bool         // whether to allow multiple sessions per user
	Client                 *RedisClient // Redis client instance to use for accessing the server
}

// SessionMeta represents metadata associated with a session.
type SessionMeta struct {
	ID      string `json:"sid"`           // ID of the session.
	UserID  int64  `json:"uid"`           // ID of the user associated with the session.
	GroupID string `json:"gid,omitempty"` // Optional ID of the group associated with the session.
}

// param returns a sessionKeyParam derived from the session metadata.
func (meta SessionMeta) param() sessionKeyParam {
	return sessionKeyParam{
		ID:      meta.ID,
		userID:  strconv.FormatInt(meta.UserID, 10),
		groupID: meta.GroupID,
	}
}

// Session represents a session with associated metadata and data.
type Session struct {
	Meta SessionMeta `json:"meta"`           // Metadata associated with the session.
	Data any         `json:"data,omitempty"` // Optional data associated with the session.
}

// Scan reads session data and decodes it into a Go value pointed
// to by dest using JSON parsing and validation.
func (session Session) Scan(dest any) error {
	b, _ := json.Marshal(session.Data)
	return JSON.ParseAndValidate(string(b), dest)
}

// sessionKeyParam represents a parameter used to uniquely identify a session.
type sessionKeyParam struct {
	ID      string // ID of the session.
	userID  string // ID of the user associated with the session.
	groupID string // Optional ID of the group associated with the session.
}

// SessionHandler represents an interface for managing sessions.
type SessionHandler interface {
	Set(session Session, expiresAt int64) error
	Get(meta SessionMeta) (Session, error)
	ListByUserID(userId int64) ([]Session, error)
	ListByGroupID(groupId string) ([]Session, error)
	Exists(meta SessionMeta) (bool, error)
	Count(uniqueByUser bool) (int, error)
	CountByUserID(userId int64, uniqueByUser bool) (int, error)
	CountByGroupID(groupId string, uniqueByUser bool) (int, error)
	Delete(meta SessionMeta) error
	DeleteByUserID(userId int64) error
	DeleteByGroupID(groupId string) error
	DeleteAll() error
}

// SessionRedisHandler is used to handle session information stored in Redis.
type SessionRedisHandler struct {
	multipleSessionPerUser bool
	prefixKey              string
	client                 *RedisClient
}

// NewSessionHandler creates a new Redis session handler using the provided configuration.
func (RedisUtil) NewSessionHandler(config SessionRedisConfig) SessionHandler {
	config.SessionKey = strings.TrimSpace(config.SessionKey)
	config.UserKey = strings.TrimSpace(config.UserKey)
	if config.SessionKey == "" {
		config.SessionKey = DefaultRedisSessionKey
	}
	if config.UserKey == "" {
		config.UserKey = DefaultRedisUserKey
	}
	var builder strings.Builder
	builder.WriteString(config.SessionKey)
	builder.WriteByte(':')
	builder.WriteString(config.UserKey)
	return &SessionRedisHandler{
		multipleSessionPerUser: config.MultipleSessionPerUser,
		prefixKey:              builder.String(),
		client:                 config.Client,
	}
}

// Set saves session data to Redis and sets an expiration time.
func (h *SessionRedisHandler) Set(s Session, expiresAt int64) error {
	key := h.getKey(s.Meta.param())
	ttl := time.Duration(Max(1, expiresAt-time.Now().Unix())) * time.Second
	return h.client.SetStruct(key, s, time.Duration(ttl.Seconds())*time.Second)
}

// Get retrieves session data from Redis using the provided metadata.
func (h *SessionRedisHandler) Get(meta SessionMeta) (Session, error) {
	key := h.getKey(meta.param())
	session := Session{}
	if err := h.client.GetStruct(key, &session); err != nil {
		if err == redis.Nil {
			return session, ErrSessionNotFound
		}
		return session, err
	}
	if session.Meta.ID != meta.ID || session.Meta.UserID != meta.UserID {
		return session, ErrSessionInvalid
	}
	return session, nil
}

// ListByUserID returns a list of sessions associated with the given user ID.
func (h *SessionRedisHandler) ListByUserID(userId int64) ([]Session, error) {
	return h.find(h.keyByUserID(userId), func(session *Session) bool {
		return session.Meta.UserID == userId
	})
}

// ListByGroupID returns a list of sessions associated with the given group ID.
// If "*" is passed as the groupId parameter, returns an empty slice and nil error.
func (h *SessionRedisHandler) ListByGroupID(groupId string) ([]Session, error) {
	if groupId == "*" {
		return []Session{}, nil
	}
	return h.find(h.keyByGroupID(groupId), func(session *Session) bool {
		return session.Meta.GroupID == groupId
	})
}

// Exists checks if session data corresponding to the given SessionMeta exists in Redis.
func (h *SessionRedisHandler) Exists(meta SessionMeta) (bool, error) {
	key := h.getKey(meta.param())
	res, err := h.client.Exists(context.TODO(), key).Result()
	return res != 0, err
}

// Count returns the total number of sessions stored in Redis
// with option to count only unique sessions.
func (h *SessionRedisHandler) Count(uniqueByUser bool) (int, error) {
	return h.countByKey(h.keyAll(), uniqueByUser)
}

// CountByUserID counts the number of sessions stored in Redis
// for a specific user ID, with option to count only unique sessions.
func (h *SessionRedisHandler) CountByUserID(userId int64, uniqueByUser bool) (int, error) {
	return h.countByKey(h.keyByUserID(userId), uniqueByUser)
}

// CountByGroupID counts the number of sessions stored in Redis for a specific group ID.
func (h *SessionRedisHandler) CountByGroupID(groupId string, uniqueByUser bool) (int, error) {
	return h.countByKey(h.keyByGroupID(groupId), uniqueByUser)
}

// Delete deletes a session from Redis based on its metadata.
func (h *SessionRedisHandler) Delete(meta SessionMeta) error {
	key := h.getKey(meta.param())
	return h.client.Del(context.TODO(), key).Err()
}

// DeleteByUserID deletes all sessions corresponding to a given user ID.
func (h *SessionRedisHandler) DeleteByUserID(userId int64) error {
	return h.deleteSessionKeys(h.keyByUserID(userId))
}

// DeleteByGroupID deletes all sessions corresponding to a given group ID.
func (h *SessionRedisHandler) DeleteByGroupID(groupId string) error {
	return h.deleteSessionKeys(h.keyByGroupID(groupId))
}

// DeleteAll deletes all sessions stored in Redis.
func (h *SessionRedisHandler) DeleteAll() error {
	return h.deleteSessionKeys(h.keyAll())
}

// keyAll returns the key that matches all session keys in Redis.
func (h *SessionRedisHandler) keyAll() string {
	return h.getKey(sessionKeyParam{
		ID:      "*",
		userID:  "*",
		groupID: "*",
	})
}

// keyByUserID returns the key that matches all session keys for a given user ID.
func (h *SessionRedisHandler) keyByUserID(userId int64) string {
	return h.getKey(sessionKeyParam{
		ID:      "*",
		userID:  strconv.FormatInt(userId, 10),
		groupID: "*",
	})
}

// keyByGroupID returns the key that matches all session keys for a given group ID.
func (h *SessionRedisHandler) keyByGroupID(groupId string) string {
	return h.getKey(sessionKeyParam{
		ID:      "*",
		userID:  "*",
		groupID: groupId,
	})
}

// getKey returns the Redis key that corresponds to a given sessionKeyParam.
func (h *SessionRedisHandler) getKey(param sessionKeyParam) string {
	var builder strings.Builder
	builder.WriteString(h.prefixKey)
	builder.WriteByte(':')
	if param.groupID != "" {
		builder.WriteString(param.groupID)
		builder.WriteByte(':')
	}
	builder.WriteString(param.userID)
	if h.multipleSessionPerUser {
		builder.WriteByte(':')
		builder.WriteString(param.ID)
	}
	return builder.String()
}

// find returns a slice of Session objects with keys matching the given pattern
// and verifying with the provided verification function.
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
			if res, err := cmd.(*redis.StringCmd).Result(); err == nil {
				var session Session
				err = json.Unmarshal([]byte(res), &session)
				if err != nil || !verifyFunc(&session) {
					continue
				}
				sessions = append(sessions, session)
			}
		}
		if cursor == 0 {
			break
		}
	}
	return sessions, nil
}

// countByKey returns the number of keys matching the given pattern.
// If uniqueByUser is true, only unique user sessions are counted.
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

// scanSessionKeys scans Redis keys matching the given pattern
// and performs the given processing function on each key.
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

// deleteSessionKeys scans Redis keys matching the given pattern and deletes them.
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
