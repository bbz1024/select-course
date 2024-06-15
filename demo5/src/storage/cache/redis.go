package cache

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/redis/go-redis/v9"
	"select-course/demo5/src/constant/config"
	"select-course/demo5/src/utils/logger"
	"time"
)

var RDB *redis.Client

func init() {
	addr := fmt.Sprintf("%s:%d", config.EnvCfg.RedisHost, config.EnvCfg.RedisPort)
	opts := &redis.Options{
		Addr:           addr,
		Password:       config.EnvCfg.RedisPwd,
		DB:             config.EnvCfg.RedisDb,
		MaxIdleConns:   config.EnvCfg.RedisMaxIdleConns,
		MaxActiveConns: config.EnvCfg.RedisMaxActiveConns,
	}
	// 测试连接,解决docker-compose容器依赖问题。
	err := retry.Do(func() error {
		RDB = redis.NewClient(opts)
		if err := RDB.Ping(context.Background()).Err(); err != nil {
			return err
		}
		return nil
	}, retry.Attempts(5), retry.Delay(time.Second))
	if err != nil {
		panic(err)
	}
	logger.Logger.Info("redis connect success")
}
