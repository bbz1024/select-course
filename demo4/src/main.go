package main

import (
	"context"
	"fmt"
	"github.com/oklog/run"
	"net/http"
	"os"
	"select-course/demo4/src/constant/config"
	"select-course/demo4/src/utils/bloom"
	"select-course/demo4/src/utils/consumer"
	"select-course/demo4/src/utils/logger"
	"select-course/demo4/src/web/router"
	"syscall"
)

func main() {
	r := router.InitApiRouter()
	Initialize()

	// groutine 编排
	g := &run.Group{}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.EnvCfg.ServerPort),
		Handler: r,
	}

	// 启动服务
	g.Add(func() error {
		if err := srv.ListenAndServe(); err != nil {
			logger.Logger.Error(err)
			return err
		}
		return nil
	}, func(err error) {
		srv.Close()
	})

	// 后台启动消费者
	g.Add(func() error {
		return consumer.SelectConsumer.Consumer()
	}, func(err error) {
		logger.Logger.Error(err)
		consumer.SelectConsumer.Close()
	})

	// 监听信号,ctrl+c
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		logger.Logger.Error(err)
		os.Exit(1)
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
