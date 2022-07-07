package redis

import (
	"context"
	"fmt"
	"gin_web/Bean"
	"gin_web/settings"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"strconv"
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

func GetUserToken(ctx context.Context, userID int64) (tokenMap map[string]string, err error) {
	formatInt := strconv.FormatInt(userID, 10) //int64 转 string
	// 按前缀扫描key
	tokenMap = make(map[string]string)
	iter := redisdb.Scan(ctx, 0, formatInt+"*", 0).Iterator()
	for iter.Next(ctx) {
		result, err := redisdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			return nil, Bean.SESSION_EXPIRED.MarkError
		}
		tokenMap[iter.Val()] = result
	}

	if err = iter.Err(); err != nil {
		return nil, Bean.SYSTEM_BUSY.MarkError //返回繁忙
	}

	return tokenMap, nil

}

func SetUserToken(ctx context.Context, key string, value string, expirationtime time.Duration) error {
	err := redisdb.Set(ctx, key, value, expirationtime).Err()
	if err != nil {
		return Bean.SYSTEM_BUSY.MarkError //返回繁忙
	}
	return nil
}
