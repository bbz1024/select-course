package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"select-course/demo2/src/constant/config"
)

var RDB *redis.Client

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:           fmt.Sprintf("%s:%s", config.EnvCfg.RedisHost, config.EnvCfg.RedisPort),
		Password:       "lzb200244",
		DB:             config.EnvCfg.RedisDb,
		MaxIdleConns:   config.EnvCfg.RedisMaxIdleConns,
		MaxActiveConns: config.EnvCfg.RedisMaxActiveConns,
	})

}
