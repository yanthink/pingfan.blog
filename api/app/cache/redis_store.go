package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strings"
	"time"
)

type RedisStore[T any] struct {
	client  *redis.Client
	options *Options
}

func NewRedisStore[T any](client *redis.Client, options ...Option) *RedisStore[T] {
	return &RedisStore[T]{
		client:  client,
		options: ApplyOptions(options...),
	}
}

func (s *RedisStore[T]) Get(ctx context.Context, key string) (result T, err error) {
	bytes, err := s.client.Get(ctx, s.buildKey(key)).Bytes()
	if err != nil {
		t := reflect.TypeOf(result)

		if t.Kind() == reflect.Map {
			result = reflect.MakeMap(t).Interface().(T)
		}

		return result, err
	}

	return s.unmarshal(bytes)
}

func (s *RedisStore[T]) Set(ctx context.Context, key string, value T, expiration time.Duration) error {
	return s.client.Set(ctx, s.buildKey(key), s.marshal(value), expiration).Err()
}

func (s *RedisStore[T]) Add(ctx context.Context, key string, value T, expiration time.Duration) error {
	return s.client.SetNX(ctx, s.buildKey(key), s.marshal(value), expiration).Err()
}

func (s *RedisStore[T]) Increment(ctx context.Context, key string, value int64) (int64, error) {
	intCmd := s.client.IncrBy(ctx, s.buildKey(key), value)

	return intCmd.Val(), intCmd.Err()
}

func (s *RedisStore[T]) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	intCmd := s.client.DecrBy(ctx, s.buildKey(key), value)

	return intCmd.Val(), intCmd.Err()
}

func (s *RedisStore[T]) Forget(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.buildKey(key)).Err()
}

func (s *RedisStore[T]) Flush(ctx context.Context) error {
	return s.client.FlushAll(ctx).Err()
}

func (s *RedisStore[T]) TTL(ctx context.Context, key string) time.Duration {
	return s.client.TTL(ctx, s.buildKey(key)).Val()
}

func (s *RedisStore[T]) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return s.client.Expire(ctx, s.buildKey(key), expiration).Err()
}

func (*RedisStore[T]) marshal(value T) []byte {
	bytes, _ := json.Marshal(value)

	return bytes
}

func (*RedisStore[T]) unmarshal(bytes []byte) (result T, err error) {
	err = json.Unmarshal(bytes, &result)

	return
}

func (s *RedisStore[T]) buildKey(key string) string {
	if !strings.HasPrefix(key, s.options.Prefix) {
		return fmt.Sprintf("%s%s", s.options.Prefix, key)
	}

	return key
}
