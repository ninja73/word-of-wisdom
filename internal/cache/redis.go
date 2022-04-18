package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type memoryCache struct {
	redisClient *redis.Client
	expiration  time.Duration
}

func NewRedisCache(redisClient *redis.Client, expiration time.Duration) *memoryCache {
	cache := &memoryCache{
		redisClient: redisClient,
		expiration:  expiration,
	}

	return cache
}

func (s *memoryCache) ContainsOrAdd(ctx context.Context, key uint64) (bool, error) {
	ok, err := s.redisClient.SetNX(ctx, strconv.FormatUint(key, 10), 1, s.expiration).Result()
	if err != nil {
		return false, err
	}
	return !ok, nil
}
