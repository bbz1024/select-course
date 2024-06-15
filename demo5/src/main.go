package main

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"select-course/demo5/src/constant/services"
	"select-course/demo5/src/utils/logger"
	"select-course/demo5/src/utils/tracing"
	"select-course/demo5/src/web/router"
	"select-course/demo5/src/web/services/courses"
	"select-course/demo5/src/web/services/users"
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

	// 启动项目
	if err := r.Run(services.WebServiceAddr); err != nil {
		logger.Logger.Info("server exit", zap.Error(err))
	}

}
