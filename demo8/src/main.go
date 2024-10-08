package main

import (
	"context"
	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"os"
	"select-course/demo8/src/constant/services"
	"select-course/demo8/src/storage/database"
	"select-course/demo8/src/utils/bloom"
	"select-course/demo8/src/utils/local"
	"select-course/demo8/src/utils/logger"
	"select-course/demo8/src/utils/tracing"
	"select-course/demo8/src/web/router"
	"select-course/demo8/src/web/services/courses"
	"select-course/demo8/src/web/services/users"
	"syscall"
)

func main() {
	// 链路
	tracer, closer := tracing.Init(services.WebServiceName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// 初始化路由
	r := router.InitApiRouter()
	// 创建grpc服务实例
	users.New()
	courses.New()
	// init mysql
	if err := database.InitMysql(); err != nil {
		logger.Logger.Error("mysql init error", zap.Error(err))
		panic(err)
		return
	}
	// init bloom
	if err := bloom.InitializeBloom(); err != nil {
		logger.Logger.Error("bloom init error", zap.Error(err))
		panic(err)
		return
	}
	//init local
	if err := local.InitLocal(); err != nil {
		logger.Logger.Error("local init error", zap.Error(err))
		panic(err)
	}
	g := run.Group{}
	g.Add(func() error {
		return r.Run(services.WebServiceAddr)
	}, func(err error) {
		logger.Logger.Error("web server exit", zap.Error(err))

	})
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	if err := g.Run(); err != nil {
		logger.Logger.Error("web server exit", zap.Error(err))
		os.Exit(1)
	}
	// 启动项目

}
