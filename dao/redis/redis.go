package redis

import (
	"context"
	"fmt"
	"gin_web/settings"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"time"
)

var redisdb *redis.Client

func Init(config *settings.RedisConfig) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5000*time.Second)
	defer cancelFunc()
	redisdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			config.Host,
			config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	_, err = redisdb.Ping(ctx).Result()
	return
}

func Close() {
	err := redisdb.Close()
	if err != nil {
		zap.L().Fatal("redis close err: ", zap.Error(err))
	}
}
