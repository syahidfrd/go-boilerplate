package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements cache operations using Redis as the backend store
type RedisCache struct {
	client *redis.Client
}

// NewRedis creates a new Redis cache implementation with the provided client
func NewRedis(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Get retrieves a value by key from Redis cache
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set stores a key-value pair in Redis cache with the specified TTL
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a key from Redis cache
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
