package redis

import "github.com/cloudsrc/api.awaymail.v1.go/src/shared/models"

type Repository interface {
	CreateInboxIndex(user models.UserSession) (err error)
	CreateSentIndex(user models.UserSession) (err error)
	SearchIndex(idx string, query ...interface{}) (res []interface{})
	AggregateIndex(idx string, query ...interface{}) (count int64)
	GetKey(key string) (doc string)
	SetKey(key string, value ...interface{})
	DelKey(key string)
	IncrKey(key string)
	GetCache(idx string, query string) (res interface{})
	SetCache(idx string, items ...interface{})
	DelCache(idx string, keys ...string) (err error)
	GetKeysPrefix(prefix string) (res []string)
	PushCache(key string, item string) (err error)
	GetPushCache(key string, items ...interface{}) (res []string, err error)
}
