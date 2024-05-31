package main

import (
	"select-course/demo1/src/constant/config"
	"select-course/demo1/src/web/router"
)

func main() {
	r := router.InitApiRouter()
	err := r.Run(":" + config.EnvCfg.ServerPort)
	if err != nil {
		panic(err)
	}
}
