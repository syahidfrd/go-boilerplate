package datastore

import (
	"github.com/go-redis/redis"
)

// NewCache will create new cache instance
func NewCache(redisURL string) *redis.Client {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	if _, err := rdb.Ping().Result(); err != nil {
		panic(err)
	}

	return rdb
}
