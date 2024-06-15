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
	"select-course/demo5/src/utils/discovery"
	"select-course/demo5/src/utils/logger"
	"select-course/demo5/src/utils/tracing"
)

func main() {
	tracer, closer := tracing.Init(services.CourseRpcServerName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	rpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_opentracing.UnaryServerInterceptor(),
		),
	)
	//注册服务
	course.RegisterCourseServiceServer(rpcServer, Course{})
	listen, err := net.Listen("tcp", services.CourseRpcServerAddr)
	if err != nil {
		logger.Logger.Error("Rpc %s listen happens error for: %v",
			zap.String("UserService", services.CourseRpcServerAddr), zap.Error(err),
		)
	}
	err = discovery.Consul.Register(
		context.Background(), discovery.Service{
			Name: services.CourseRpcServerName,
			Port: services.CourseRpcServerAddr,
		},
	)

	g := &run.Group{}
	g.Add(func() error {
		return rpcServer.Serve(listen)
	}, func(err error) {
		logger.LogServer(services.CourseRpcServerName).Error("server stopped", zap.Error(err))
		rpcServer.GracefulStop()
	})
	if err := g.Run(); err != nil {
		logger.LogServer(services.CourseRpcServerName).Error("run error", zap.Error(err))
	}
}
