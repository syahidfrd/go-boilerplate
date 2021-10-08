package redis_test

import (
	"log"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	redisRepo "github.com/syahidfrd/go-boilerplate/repository/redis"
)

func SetupRedis() redisRepo.RedisRepository {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub redis connection", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	redisRepository := redisRepo.NewRedisRepository(client)
	return redisRepository
}

func TestSet(t *testing.T) {
	redisRepository := SetupRedis()
	err := redisRepository.Set("ping", "pong", time.Duration(0))
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	redisRepository := SetupRedis()
	key, val, exp := "ping", "pong", time.Duration(0)

	value, err := redisRepository.Get(key)
	assert.NotNil(t, err)
	assert.Equal(t, value, "")

	err = redisRepository.Set(key, val, exp)
	assert.NoError(t, err)

	value, err = redisRepository.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, val)
}
