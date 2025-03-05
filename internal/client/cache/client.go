package cache

import (
	"context"
	"time"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}) error
	HSet(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	HGetAll(ctx context.Context, key string) ([]interface{}, error)
	Del(ctx context.Context, key string) error
	HDel(ctx context.Context, key string, fields ...string) error
	Expire(ctx context.Context, key string, duration time.Duration) error
	Ping(ctx context.Context) error
}
