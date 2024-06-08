package main

import (
	"select-course/demo4/src/constant/config"
	"select-course/demo4/src/utils/bloom"
	"select-course/demo4/src/utils/consumer"
	"select-course/demo4/src/utils/logger"
	"select-course/demo4/src/web/router"
)

func main() {
	r := router.InitApiRouter()
	Initialize()
	err := r.Run(":" + config.EnvCfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
func Initialize() {
	var err error
	// 初始化布隆过滤器
	err = bloom.InitializeBloom()
	if err != nil {
		panic(err)
		return
	}
	// 初始化队列
	if err := consumer.InitSelectListener(); err != nil {
		logger.Logger.Error(err)
		panic(err)
	}

}
