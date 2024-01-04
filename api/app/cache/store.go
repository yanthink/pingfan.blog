package cache

import (
	"context"
	"time"
)

type Store[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T, expiration time.Duration) error
	Add(ctx context.Context, key string, value T, expiration time.Duration) error
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string, value int64) (int64, error)
	Forget(ctx context.Context, key string) error
	Flush(ctx context.Context) error
	TTL(ctx context.Context, key string) time.Duration
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
