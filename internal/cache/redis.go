package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

// redisCache implements cache operations using Redis as the backend store
type redisCache struct {
	client *redis.Client
}

// NewRedis creates a new Redis cache implementation with the provided client
func NewRedis(client *redis.Client) *redisCache {
	return &redisCache{
		client: client,
	}
}

// Get retrieves a value by key from Redis cache
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(key).Result()
}

// Set stores a key-value pair in Redis cache with the specified TTL
func (r *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(key, value, ttl).Err()
}

// Delete removes a key from Redis cache
func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(key).Err()
}
