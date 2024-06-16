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
	"select-course/demo5/src/rpc/user"
	"select-course/demo5/src/storage/database"
	"select-course/demo5/src/utils/discovery"
	"select-course/demo5/src/utils/local"
	"select-course/demo5/src/utils/logger"
	"select-course/demo5/src/utils/tracing"
)

func main() {
	// tracing
	tracer, closer := tracing.Init(services.UserRpcServerName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	rpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_opentracing.UnaryServerInterceptor(),
		),
	)
	// init server
	userServer := &User{}
	user.RegisterUserServiceServer(rpcServer, userServer)
	listen, err := net.Listen("tcp", services.UserRpcServerAddr)
	if err != nil {
		logger.Logger.Error("Rpc %s listen happens error for: %v",
			zap.String("UserService", services.UserRpcServerAddr), zap.Error(err),
		)
		panic(err)
	}
	err = discovery.Consul.Register(
		context.Background(), discovery.Service{
			Name: services.UserRpcServerName,
			Port: services.UserRpcServerAddr,
		},
	)

	// init instance
	userServer.New()
	if err := database.InitMysql(); err != nil {
		logger.Logger.Error("mysql init error", zap.Error(err))
		panic(err)
	}
	//init local
	if err := local.InitLocal(); err != nil {
		logger.Logger.Error("local init error", zap.Error(err))
		panic(err)
	}
	g := &run.Group{}
	g.Add(func() error {
		return rpcServer.Serve(listen)
	}, func(err error) {
		logger.Logger.Error("Rpc %s listen happens error for: %v",
			zap.String("UserService", services.UserRpcServerAddr), zap.Error(err),
		)
		rpcServer.GracefulStop()
	})
	if err := g.Run(); err != nil {
		logger.Logger.Error("Rpc %s listen happens error for: %v",
			zap.String("UserService", services.UserRpcServerAddr), zap.Error(err),
		)

		panic(err)
	}
}
