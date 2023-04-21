package redis

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudsrc/api.awaymail.v1.go/libs"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"
	"github.com/garyburd/redigo/redis"
)

type inboxRepository struct {
	redis *redis.Pool
}

func NewRedisInbox(redis *libs.RedisClient) *inboxRepository {
	if redis == nil {
		panic("redis client is nil")
	}
	return &inboxRepository{
		redis: redis.Redis,
	}
}

func (i *inboxRepository) CreateInboxIndex(user models.UserSession) (err error) {
	idx := fmt.Sprintf("user:inbox:%s", user.UserID.Hex())
	schema := fmt.Sprintf(`ON hash PREFIX 1 user:inbox:%s: SCHEMA `+
		`body_text TEXT NOSTEM `+
		`from TEXT NOSTEM `+
		`subject TEXT NOSTEM `+
		`received_at NUMERIC SORTABLE`,
		user.UserID.Hex())
	args := redis.Args{idx}
	for _, s := range strings.Split(schema, " ") {
		args = append(args, s)
	}
	rConn := i.redis.Get()
	defer rConn.Close()
	_, err = rConn.Do("FT.CREATE", args...)
	return
}

func (i *inboxRepository) CreateSentIndex(user models.UserSession) (err error) {
	idx := fmt.Sprintf("user:sent:%s", user.UserID.Hex())
	schema := fmt.Sprintf(`ON hash PREFIX 1 user:sent:%s: SCHEMA `+
		`body_text TEXT NOSTEM `+
		`from TEXT NOSTEM `+
		`subject TEXT NOSTEM `+
		`sent NUMERIC SORTABLE`,
		user.UserID.Hex())
	args := redis.Args{idx}
	for _, s := range strings.Split(schema, " ") {
		args = append(args, s)
	}
	rConn := i.redis.Get()
	defer rConn.Close()
	_, err = rConn.Do("FT.CREATE", args...)
	return
}

func (i *inboxRepository) SearchIndex(idx string, query ...interface{}) (res []interface{}) {
	args := redis.Args{idx}
	args = append(args, query...)
	rConn := i.redis.Get()
	defer rConn.Close()
	resp, _ := redis.Values(rConn.Do("FT.SEARCH", args...))
	if resp == nil {
		return nil
	}
	for i := 0; i <= len(resp)-1; i++ {
		var doc interface{}
		switch resp[i].(type) {
		case []interface{}:
			json.Unmarshal(resp[i].([]interface{})[1].([]byte), &doc)
			res = append(res, doc)
		}
	}
	return res
}

func (i *inboxRepository) AggregateIndex(idx string, query ...interface{}) (count int64) {
	args := redis.Args{idx}
	args = append(args, query...)
	rConn := i.redis.Get()
	defer rConn.Close()
	resp, _ := redis.Values(rConn.Do("FT.AGGREGATE", args...))
	if resp == nil {
		return
	}
	for i := 0; i <= len(resp)-1; i++ {
		switch resp[i].(type) {
		case []interface{}:
			count++
		}
	}
	return
}

func (i *inboxRepository) GetKey(key string) (doc string) {
	rConn := i.redis.Get()
	defer rConn.Close()
	res, _ := rConn.Do("GET", key)
	if res == nil {
		return ""
	}
	return string(res.([]byte))
}

func (i *inboxRepository) SetKey(key string, value ...interface{}) {
	rConn := i.redis.Get()
	args := redis.Args{key}
	args = append(args, value...)
	defer rConn.Close()
	_, _ = rConn.Do("SET", args...)
}

func (i *inboxRepository) DelKey(key string) {
	rConn := i.redis.Get()
	args := redis.Args{key}
	defer rConn.Close()
	_, _ = rConn.Do("DEL", args...)
}

func (i *inboxRepository) IncrKey(key string) {
	rConn := i.redis.Get()
	defer rConn.Close()
	_, _ = rConn.Do("INCR", key)
}

func (i *inboxRepository) GetCache(idx string, query string) (res interface{}) {
	rConn := i.redis.Get()
	args := redis.Args{idx}
	args = append(args, query)
	defer rConn.Close()
	val, _ := rConn.Do("HGET", args...)
	if val == nil {
		return nil
	}
	json.Unmarshal(val.([]byte), &res)
	return res
}

func (i *inboxRepository) SetCache(idx string, items ...interface{}) {
	args := redis.Args{idx}
	args = append(args, items...)
	rConn := i.redis.Get()
	defer rConn.Close()
	rConn.Do("HSET", args...)
}

func (i *inboxRepository) DelCache(idx string, keys ...string) (err error) {
	args := redis.Args{idx}
	for _, s := range keys {
		args = append(args, s)
	}
	rConn := i.redis.Get()
	defer rConn.Close()
	if _, err := rConn.Do("HDEL", args...); err != nil {
		return err
	}
	return nil
}

func (i *inboxRepository) GetKeysPrefix(prefix string) (res []string) {
	rConn := i.redis.Get()
	defer rConn.Close()
	val, _ := redis.Values(rConn.Do("KEYS", prefix))
	for _, v := range val {
		res = append(res, string(v.([]byte)))
	}
	return res
}

func (i *inboxRepository) PushCache(key string, item string) (err error) {
	rConn := i.redis.Get()
	args := redis.Args{key}
	args = append(args, item)
	defer rConn.Close()
	_, err = rConn.Do("LPUSH", args...)
	return
}

func (i *inboxRepository) GetPushCache(key string, items ...interface{}) (res []string, err error) {
	rConn := i.redis.Get()
	args := redis.Args{key}
	args = append(args, items...)
	defer rConn.Close()
	val, err := redis.Values(rConn.Do("LRANGE", args...))
	for _, v := range val {
		res = append(res, string(v.([]byte)))
	}
	return
}
