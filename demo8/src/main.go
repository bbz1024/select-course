package main

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"select-course/demo8/src/constant/services"
	"select-course/demo8/src/storage/database"
	"select-course/demo8/src/utils/bloom"
	"select-course/demo8/src/utils/local"
	"select-course/demo8/src/utils/logger"
	"select-course/demo8/src/utils/tracing"
	"select-course/demo8/src/web/router"
	"select-course/demo8/src/web/services/courses"
	"select-course/demo8/src/web/services/users"
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
	// 启动项目
	if err := r.Run(services.WebServiceAddr); err != nil {
		logger.Logger.Info("server exit", zap.Error(err))
		panic(err)
	}

}
