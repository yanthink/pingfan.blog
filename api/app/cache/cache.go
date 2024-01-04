package cache

import (
	"blog/app"
	"blog/config"
	"context"
	"time"
)

type Cache[T any] struct {
	Store[T]
}

func New[T any]() *Cache[T] {
	switch config.Cache.Store {
	default:
		return &Cache[T]{
			Store: NewRedisStore[T](app.Redis, WithPrefix(config.Redis.Prefix)),
		}
	}
}

func (c *Cache[T]) Pull(ctx context.Context, key string) (T, error) {
	value, err := c.Get(ctx, key)
	_ = c.Forget(ctx, key)

	return value, err
}

func (c *Cache[T]) Put(ctx context.Context, key string, value T, expiration time.Duration) error {
	return c.Set(ctx, key, value, expiration)
}

func (c *Cache[T]) Remember(ctx context.Context, key string, expiration time.Duration, callback func() T) (result T, err error) {
	if result, err = c.Get(ctx, key); err == nil {
		return
	}

	result = callback()
	err = c.Put(ctx, key, result, expiration)

	return
}
