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

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionInvalid  = errors.New("invalid session")
)

// Default redis session constants
const (
	DefaultRedisSessionKey = "session"
	DefaultRedisUserKey    = "user"
)

const (
	sessionRedisLimitTx = 100
)

// SessionRedisConfig represents the configuration used to setup a Redis session manager.
type SessionRedisConfig struct {
	SessionKey             string       // The key used to store session data in Redis
	UserKey                string       // The key used to identify the user associated with a given session
	MultipleSessionPerUser bool         // Indicates whether multiple sessions are allowed per user
	Client                 *RedisClient // The Redis client used to communicate with the Redis server
}

type SessionRedisGetKeyParam struct {
	SessionID string
	UserID    string
	GroupID   string
}

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
	Count(uniqueUser bool) (int, error)
	CountByUserID(userId int64, uniqueUser bool) (int, error)
	CountByGroupID(groupId string, uniqueUser bool) (int, error)
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
	SessionHandler         // TODO! remove
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

// GetKey returns the Redis key for the given session.
func (h *SessionRedisHandler) GetKey(param SessionRedisGetKeyParam) string {
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

// Set stores the provided session in Redis.
func (h *SessionRedisHandler) Set(s Session, expiresAt int64) error {
	key := h.GetKey(SessionRedisGetKeyParam{
		SessionID: s.ID,
		UserID:    strconv.FormatInt(s.UserID, 10),
		GroupID:   s.GroupID,
	})
	ttl := time.Duration(Max(1, expiresAt-time.Now().Unix())) * time.Second
	return h.client.SetStruct(key, s, time.Duration(ttl)*time.Second)
}

// Get retrieves the session with the specified ID for the given user and group.
func (h *SessionRedisHandler) Get(sessionId string, userId int64, GroupId string) (Session, error) {
	key := h.GetKey(SessionRedisGetKeyParam{
		SessionID: sessionId,
		UserID:    strconv.FormatInt(userId, 10),
		GroupID:   GroupId,
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

func (h *SessionRedisHandler) Find(key string, verifyFunc func(*Session) bool) ([]Session, error) {
	var cursor uint64
	ctx := context.TODO()
	sessions := []Session{}
	for {
		var err error
		var keys []string
		keys, cursor, err = h.client.Scan(ctx, cursor, key, sessionRedisLimitTx).Result()
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

// ListByUserID retrieves all sessions for the given user.
func (h *SessionRedisHandler) ListByUserID(userId int64) ([]Session, error) {
	key := h.GetKey(SessionRedisGetKeyParam{
		SessionID: "*",
		UserID:    strconv.FormatInt(userId, 10),
		GroupID:   "*",
	})
	return h.Find(key, func(s *Session) bool {
		return s.UserID == userId
	})
}

// ListByGroupID retrieves all sessions for the given group.
func (h *SessionRedisHandler) ListByGroupID(groupId string) ([]Session, error) {
	if groupId == "*" {
		return []Session{}, nil
	}
	key := h.GetKey(SessionRedisGetKeyParam{
		SessionID: "*",
		UserID:    "*",
		GroupID:   groupId,
	})
	return h.Find(key, func(s *Session) bool {
		return s.GroupID == groupId
	})
}

// // Exists checks if the session with the specified ID for the given user and group exists.
// func (handler *SessionRedisHandler) Exists(sid string, uid int64, gid string) (bool, error) {
// 	key := fmt.Sprintf("%s:%d:%s", handler.prefixKey, uid, sid)
// 	res, err := handler.client.Exists(key).Result()
// 	if err != nil {
// 		return false, err
// 	}
// 	if res == 0 {
// 		return false, nil
// 	}
// 	if !handler.multipleSessionPerUser && gid != "" {
// 		// check if the session belongs to the given group
// 		session, err := handler.Get(sid, uid, gid)
// 		if err != nil {
// 			return false, err
// 		}
// 		if session.GroupID != gid {
// 			return false, nil
// 		}
// 	}
// 	return true, nil
// }

// // Count returns the number of stored sessions.
// func (handler *SessionRedisHandler) Count(bool) (int, error) {
// 	key := fmt.Sprintf("%s:*", handler.prefixKey)
// 	keys, err := handler.client.Keys(key).Result()
// 	if err != nil {
// 		return -1, err
// 	}
// 	return len(keys), nil
// }

// // CountByUserID returns the number of sessions associated with the given user.
// func (handler *SessionRedisHandler) CountByUserID(uid int64, uniqueUser bool) (int, error) {
// 	key := fmt.Sprintf("%s:%d:*", handler.prefixKey, uid)
// 	keys, err := handler.client.Keys(key).Result()
// 	if err != nil {
// 		return -1, err
// 	}
// 	return len(keys), nil
// }

// // CountByGroupID returns the number of sessions associated with the given group.
// func (handler *SessionRedisHandler) CountByGroupID(gid string, uniqueUser bool) (int, error) {
// 	key := fmt.Sprintf("%s:*", handler.prefixKey)
// 	keys, err := handler.client.Keys(key).Result()
// 	if err != nil {
// 		return -1, err
// 	}
// 	count := 0
// 	for _, k := range keys {
// 		res, err := handler.client.Get(k).Result()
// 		if err == redis.Nil {
// 			continue
// 		} else if err != nil {
// 			return -1, err
// 		}
// 		var session Session
// 		err = json.Unmarshal([]byte(res), &session)
// 		if err != nil {
// 			continue
// 		}
// 		if session.GroupID == gid {
// 			count++
// 		}
// 	}
// 	return count, nil
// }

// // Delete removes the session with the specified ID for the given user and group from Redis.
// func (handler *SessionRedisHandler) Delete(sid string, uid int64, gid string) error {
// 	if !handler.multipleSessionPerUser && gid != "" {
// 		// check if the session belongs to the given group
// 		session, err := handler.Get(sid, uid, gid)
// 		if err != nil {
// 			return err
// 		}
// 		if session.GroupID != gid {
// 			return fmt.Errorf("Session not found")
// 		}
// 	}
// 	key := fmt.Sprintf("%s:%d:%s", handler.prefixKey, uid, sid)
// 	err := handler.client.Del(key).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // DeleteByUserID removes all sessions associated with the given user from Redis.
// func (handler *SessionRedisHandler) DeleteByUserID(uid int64) error {
// 	key := fmt.Sprintf("%s:%d:*", handler.prefixKey, uid)
// 	keys, err := handler.client.Keys(key).Result()
// 	if err != nil {
// 		return err
// 	}
// 	for _, k := range keys {
// 		err = handler.client.Del(k).Err()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // DeleteByGroupID removes all sessions associated with the given group from Redis.
// func (handler *SessionRedisHandler) DeleteByGroupID(gid string) error {
// 	key := fmt.Sprintf("%s:*", handler.prefixKey)
// 	keys, err := handler.client.Keys(key).Result()
// 	if err != nil {
// 		return err
// 	}
// 	for _, k := range keys {
// 		res, err := handler.client.Get(k).Result()
// 		if err == redis.Nil {
// 			continue
// 		} else if err != nil {
// 			return err
// 		}
// 		var session Session
// 		err = json.Unmarshal([]byte(res), &session)
// 		if err != nil {
// 			continue
// 		}
// 		if session.GroupID == gid {
// 			err = handler.client.Del(k).Err()
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// // DeleteAll removes all sessions from Redis.
// func (handler *SessionRedisHandler) DeleteAll() error {
// 	key := fmt.Sprintf("%s:*", handler.prefixKey)
// 	keys, err := handler.client.Keys(key).Result()
// 	if err != nil {
// 		return err
// 	}
// 	for _, k := range keys {
// 		err = handler.client.Del(k).Err()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
