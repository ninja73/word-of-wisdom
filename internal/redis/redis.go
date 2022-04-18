package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func InitRedis(host string, port int, password string, db, poolSize int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
