package models

import (
	"blog/app"
	"blog/config"
	"context"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type RedisIncrField struct {
	Table string
	Field string
	Date  time.Time
}

func NewRedisIncrField(table string, field string, date time.Time) *RedisIncrField {
	return &RedisIncrField{
		Table: table,
		Field: field,
		Date:  date,
	}
}

func (field *RedisIncrField) Current(id any) (current int64) {
	ctx := context.Background()
	hashKey := field.getHashKey()
	hashField := field.getHashField(id)

	current, _ = app.Redis.HGet(ctx, hashKey, hashField).Int64()

	return
}

func (field *RedisIncrField) Incr(id any, value int64) int64 {
	ctx := context.Background()
	hashKey := field.getHashKey()
	hashField := field.getHashField(id)

	return app.Redis.HIncrBy(ctx, hashKey, hashField, value).Val()
}

func (field *RedisIncrField) Del(id any) int64 {
	ctx := context.Background()
	hashKey := field.getHashKey()
	hashField := field.getHashField(id)

	return app.Redis.HDel(ctx, hashKey, hashField).Val()
}

func (field *RedisIncrField) getHashKey() string {
	return fmt.Sprintf("%sincr_field:%s:%s", config.Redis.Prefix, field.Table, field.Date.Format(time.DateOnly))
}

func (field *RedisIncrField) getHashField(id any) string {
	return fmt.Sprintf("%s_%v", field.Field, id)
}

func (field *RedisIncrField) Scan(fn func(id, value string)) {
	ctx := context.Background()
	hashKey := field.getHashKey()

	iter := app.Redis.HScan(ctx, hashKey, 0, "", 10).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()

		if iter.Next(ctx) {
			id := strings.Replace(key, fmt.Sprintf("%s_", field.Field), "", 1)
			fn(id, iter.Val())
		}
	}
}

func (field *RedisIncrField) SyncToDatabase() {
	field.Scan(func(id, value string) {
		app.DB.
			Table(field.Table).
			Where("id = ?", id).
			UpdateColumn(field.Field, gorm.Expr(fmt.Sprintf("%s + %s", field.Field, value)))
	})

	field.Flush()
}

func (field *RedisIncrField) Flush() {
	ctx := context.Background()
	hashKey := field.getHashKey()

	app.Redis.Del(ctx, hashKey)
}
