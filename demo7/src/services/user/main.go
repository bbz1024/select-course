package main

import (
	"context"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"select-course/demo7/src/constant/services"
	"select-course/demo7/src/rpc/user"
	"select-course/demo7/src/storage/database"
	"select-course/demo7/src/utils/discovery"
	"select-course/demo7/src/utils/local"
	"select-course/demo7/src/utils/logger"
	"select-course/demo7/src/utils/tracing"
	"syscall"
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

	// init instance
	userServer.New()
	if err = discovery.Consul.Register(
		context.Background(), discovery.Service{
			Name: services.UserRpcServerName,
			Port: services.UserRpcServerAddr,
		},
	); err != nil {
		panic(err)
	}

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
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	if err := g.Run(); err != nil {
		logger.Logger.Error("Rpc %s listen happens error for: %v",
			zap.String("UserService", services.UserRpcServerAddr), zap.Error(err),
		)
		if err := discovery.Consul.Deregister(context.Background(), services.UserRpcServerAddr); err != nil {
			logger.LogService(services.UserRpcServerAddr).Error("deregister error", zap.Error(err))
		}
		os.Exit(1)
	}
}
