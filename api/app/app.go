package app

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	Redis   *redis.Client
	Redsync *redsync.Redsync
	Limiter *redis_rate.Limiter
	Logger  *zap.Logger
)
