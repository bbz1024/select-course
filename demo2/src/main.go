package main

import (
	"select-course/demo2/src/constant/config"
	"select-course/demo2/src/utils/bloom"
	"select-course/demo2/src/web/router"
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
	// 初始化布隆过滤器
	err := bloom.InitializeBloom()
	if err != nil {
		panic(err)
		return
	}
}
