package health

import (
	"context"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// store implements health check data operations using GORM and Redis
type store struct {
	dbConn      *gorm.DB
	redisClient *redis.Client
}

// NewStore creates a new health store with the provided database and Redis connections
func NewStore(dbConn *gorm.DB, redisClient *redis.Client) *store {
	return &store{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

// PingDatabase checks database connectivity
func (s *store) PingDatabase(ctx context.Context) error {
	sqlDB, err := s.dbConn.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// PingCache checks Redis cache connectivity
func (s *store) PingCache(ctx context.Context) error {
	return s.redisClient.Ping().Err()
}
