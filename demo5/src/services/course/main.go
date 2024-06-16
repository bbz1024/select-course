package main

import (
	"context"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"select-course/demo5/src/constant/services"
	"select-course/demo5/src/rpc/course"
	"select-course/demo5/src/storage/database"
	"select-course/demo5/src/utils/consumer"
	"select-course/demo5/src/utils/discovery"
	"select-course/demo5/src/utils/local"
	"select-course/demo5/src/utils/logger"
	"select-course/demo5/src/utils/tracing"
)

func main() {
	// tracing init
	tracer, closer := tracing.Init(services.CourseRpcServerName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// rpc init
	rpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_opentracing.UnaryServerInterceptor(),
		),
	)
	courseService := &Course{}
	course.RegisterCourseServiceServer(rpcServer, courseService)

	// init instance
	courseService.New()
	listen, err := net.Listen("tcp", services.CourseRpcServerAddr)
	if err != nil {
		logger.Logger.Error("Rpc %s listen happens error for: %v",
			zap.String("UserService", services.CourseRpcServerAddr), zap.Error(err),
		)
		panic(err)
	}
	// init mysql
	if err := database.InitMysql(); err != nil {
		logger.Logger.Error("mysql init error", zap.Error(err))
		panic(err)
	}
	//init local
	if err := local.InitLocal();err!=nil{
		logger.Logger.Error("local init error", zap.Error(err))
		panic(err)
	}
	// init mq
	if err := consumer.InitSelectListener(); err != nil {
		logger.Logger.Error("SelectConsumer init error for: %v", zap.Error(err))
		panic(err)

	}

	// consul register
	err = discovery.Consul.Register(
		context.Background(), discovery.Service{
			Name: services.CourseRpcServerName,
			Port: services.CourseRpcServerAddr,
		},
	)

	// -------------------- g-routine compose  --------------------
	g := &run.Group{}
	g.Add(func() error {
		return rpcServer.Serve(listen)
	}, func(err error) {
		logger.LogService(services.CourseRpcServerName).Error("server stopped", zap.Error(err))
		rpcServer.GracefulStop()
	})

	// back consumer
	g.Add(func() error {
		return consumer.SelectConsumer.Consumer()
	}, func(err error) {
		logger.LogService(services.CourseRpcServerName).Error("consumer stopped", zap.Error(err))
		consumer.SelectConsumer.Close()
	})
	if err := g.Run(); err != nil {
		logger.LogService(services.CourseRpcServerName).Error("run error", zap.Error(err))
		panic(err)
	}
}
