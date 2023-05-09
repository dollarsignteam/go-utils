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

type SessionRedisConfig struct {
	SessionKey             string
	UserKey                string
	MultipleSessionPerUser bool
	Client                 *RedisClient
}

type SessionMeta struct {
	ID      string `json:"sid"`
	UserID  int64  `json:"uid"`
	GroupID string `json:"gid,omitempty"`
}

func (meta SessionMeta) param() sessionKeyParam {
	return sessionKeyParam{
		ID:      meta.ID,
		userID:  strconv.FormatInt(meta.UserID, 10),
		groupID: meta.GroupID,
	}
}

type Session struct {
	Meta SessionMeta `json:"meta"`
	Data any         `json:"data,omitempty"`
}

type sessionKeyParam struct {
	ID      string
	userID  string
	groupID string
}

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

type SessionRedisHandler struct {
	multipleSessionPerUser bool
	prefixKey              string
	client                 *RedisClient
}

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

func (h *SessionRedisHandler) Set(s Session, expiresAt int64) error {
	key := h.getKey(s.Meta.param())
	ttl := time.Duration(Max(1, expiresAt-time.Now().Unix())) * time.Second
	return h.client.SetStruct(key, s, time.Duration(ttl)*time.Second)
}

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

func (h *SessionRedisHandler) ListByUserID(userId int64) ([]Session, error) {
	return h.find(h.keyByUserID(userId), func(session *Session) bool {
		return session.Meta.UserID == userId
	})
}

func (h *SessionRedisHandler) ListByGroupID(groupId string) ([]Session, error) {
	if groupId == "*" {
		return []Session{}, nil
	}
	return h.find(h.keyByGroupID(groupId), func(session *Session) bool {
		return session.Meta.GroupID == groupId
	})
}

func (h *SessionRedisHandler) Exists(meta SessionMeta) (bool, error) {
	key := h.getKey(meta.param())
	res, err := h.client.Exists(context.TODO(), key).Result()
	return res != 0, err
}

func (h *SessionRedisHandler) Count(uniqueByUser bool) (int, error) {
	return h.countByKey(h.keyAll(), uniqueByUser)
}

func (h *SessionRedisHandler) CountByUserID(userId int64, uniqueByUser bool) (int, error) {
	return h.countByKey(h.keyByUserID(userId), uniqueByUser)
}

func (h *SessionRedisHandler) CountByGroupID(groupId string, uniqueByUser bool) (int, error) {
	return h.countByKey(h.keyByGroupID(groupId), uniqueByUser)
}

func (h *SessionRedisHandler) Delete(meta SessionMeta) error {
	key := h.getKey(meta.param())
	return h.client.Del(context.TODO(), key).Err()
}

func (h *SessionRedisHandler) DeleteByUserID(userId int64) error {
	return h.deleteSessionKeys(h.keyByUserID(userId))
}

func (h *SessionRedisHandler) DeleteByGroupID(groupId string) error {
	return h.deleteSessionKeys(h.keyByGroupID(groupId))
}

func (h *SessionRedisHandler) DeleteAll() error {
	return h.deleteSessionKeys(h.keyAll())
}

func (h *SessionRedisHandler) keyAll() string {
	return h.getKey(sessionKeyParam{
		ID:      "*",
		userID:  "*",
		groupID: "*",
	})
}

func (h *SessionRedisHandler) keyByUserID(userId int64) string {
	return h.getKey(sessionKeyParam{
		ID:      "*",
		userID:  strconv.FormatInt(userId, 10),
		groupID: "*",
	})
}

func (h *SessionRedisHandler) keyByGroupID(groupId string) string {
	return h.getKey(sessionKeyParam{
		ID:      "*",
		userID:  "*",
		groupID: groupId,
	})
}

func (h *SessionRedisHandler) getKey(param sessionKeyParam) string {
	var builder strings.Builder
	builder.WriteString(h.prefixKey)
	builder.WriteByte(':')
	builder.WriteString(param.groupID)
	builder.WriteByte(':')
	builder.WriteString(param.userID)
	if h.multipleSessionPerUser {
		builder.WriteByte(':')
		builder.WriteString(param.ID)
	}
	return builder.String()
}

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
			if err != nil || !verifyFunc(&session) {
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
