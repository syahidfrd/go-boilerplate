package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// RedisRepository represent the redis repositories
type RedisRepository interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
}

type redisRepository struct {
	client *redis.Client
}

// NewRedisRepository will create an object that represent the RedisRepository interface
func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{
		client: client,
	}
}

// Set attaches the redis repository and set the data
func (r *redisRepository) Set(key string, value interface{}, exp time.Duration) error {
	return r.client.Set(key, value, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *redisRepository) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}
