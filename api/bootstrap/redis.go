package bootstrap

import (
	"blog/app"
	"blog/config"
	"context"
	"fmt"
	"github.com/go-redis/redis_rate/v10"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"net"
)

type RedisHook struct{}

func (RedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		if conn, err = next(ctx, network, addr); err != nil {
			app.Logger.Sugar().Error(err)
		}

		return
	}
}

func (RedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		return next(ctx, cmd)
	}
}

func (RedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

func SetupRedis() {
	app.Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Pass,
		DB:       config.Redis.DB,
	})

	if err := app.Redis.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	app.Redis.AddHook(RedisHook{})

	app.Redsync = redsync.New(goredis.NewPool(app.Redis))
	app.Limiter = redis_rate.NewLimiter(app.Redis)
}
