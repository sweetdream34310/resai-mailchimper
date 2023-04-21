package libs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/garyburd/redigo/redis"
)

type Document struct {
	Id         string
	Payload    []byte
	Properties map[string]interface{}
}

type RedisClient struct {
	Redis *redis.Pool
	dbno  int
}

func newRedisClient(config config.Config) *RedisClient {

	pool := &redis.Pool{
		MaxIdle:   config.Redis.MaxIdle,
		MaxActive: config.Redis.MaxActive,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
				[]redis.DialOption{
					redis.DialPassword(config.Redis.Password),
				}...,
			)
		},
	}
	pool.TestOnBorrow = func(c redis.Conn, t time.Time) (err error) {
		if time.Since(t) > time.Second {
			_, err = c.Do("PING")
		}
		return err
	}
	return &RedisClient{
		Redis: pool,
		dbno:  config.Redis.DB,
	}
}

// Close : close the session to redis
func (r *RedisClient) Close() (err error) {
	return r.Redis.Close()
}

// ReadCache : Read from redis cache and returns the object.
func (r *RedisClient) GetCache(idx string, query string) (res interface{}) {
	args := redis.Args{idx}
	args = append(args, query)
	rConn := r.Redis.Get()
	defer rConn.Close()
	val, _ := rConn.Do("HGET", args...)
	if val == nil {
		return nil
	}
	json.Unmarshal(val.([]byte), &res)
	return res
}

// SetKey : Sets a redis key
func (r *RedisClient) SetKey(key string, value ...interface{}) {
	args := redis.Args{key}
	args = append(args, value...)
	rConn := r.Redis.Get()
	defer rConn.Close()
	rConn.Do("SET", args...)
}

// Getkey : Gets a redis key
func (r *RedisClient) Getkey(key string) (doc string) {
	rConn := r.Redis.Get()
	defer rConn.Close()
	res, _ := rConn.Do("GET", key)
	if res == nil {
		return ""
	}
	return string(res.([]byte))
}

// Incrkey : Increments the value of a redis key
func (r *RedisClient) Incrkey(key string) {
	rConn := r.Redis.Get()
	defer rConn.Close()
	rConn.Do("INCR", key)
}

// ReadCache : Read from redis cache and returns the object.
func (r *RedisClient) SearchIndex(idx string, query ...interface{}) (res []interface{}) {
	args := redis.Args{idx}
	args = append(args, query...)
	rConn := r.Redis.Get()
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

// SetCache : Set to redis cache with expiry.
func (r *RedisClient) SetCache(idx string, items ...interface{}) {
	args := redis.Args{idx}
	args = append(args, items...)
	rConn := r.Redis.Get()
	defer rConn.Close()
	rConn.Do("HSET", args...)
}

func (r *RedisClient) CreateIndex(idx string, schema string) (err error) {
	args := redis.Args{idx}
	for _, s := range strings.Split(schema, " ") {
		args = append(args, s)
	}
	rConn := r.Redis.Get()
	defer rConn.Close()
	if _, err := rConn.Do("FT.CREATE", args...); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) DelCache(idx string, keys ...string) (err error) {
	args := redis.Args{idx}
	for _, s := range keys {
		args = append(args, s)
	}
	rConn := r.Redis.Get()
	defer rConn.Close()
	if _, err := rConn.Do("HDEL", args...); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) GetKeysPrefix(prefix string) (res []string) {
	rConn := r.Redis.Get()
	defer rConn.Close()
	val, _ := redis.Values(rConn.Do("KEYS", prefix))
	for _, v := range val {
		res = append(res, string(v.([]byte)))
	}
	return res
}
