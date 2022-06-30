package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var redisdb *redis.Client

func Init() (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5000*time.Second)
	defer cancelFunc()
	redisdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
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
