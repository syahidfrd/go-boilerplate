package datastore

import (
	"github.com/go-redis/redis"
)

// NewCache will create new cache instance
func NewCache(redisURL string) (client *redis.Client, err error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return
	}

	client = redis.NewClient(opt)
	if _, err = client.Ping().Result(); err != nil {
		return
	}

	return
}
