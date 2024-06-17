package grpc

import (
	"context"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"select-course/demo7/src/utils/discovery"
	"select-course/demo7/src/utils/logger"
	"time"
)

func Dial(addr string) (*grpc.ClientConn, error) {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: false,            // send pings even without active streams
	}
	return grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithChainUnaryInterceptor(
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				ctx = opentracing.ContextWithSpan(ctx, opentracing.SpanFromContext(ctx))
				return invoker(ctx, method, req, reply, cc, opts...)
			},
			otelgrpc.UnaryClientInterceptor(),
			//注入trace
			grpc_opentracing.UnaryClientInterceptor(),
		),
	)
}

func Connect(ctx context.Context, serviceName string) *grpc.ClientConn {
	addr, err := discovery.Consul.GetService(ctx, serviceName)
	if err != nil {
		logger.LogService(serviceName).Error("get service error", zap.Error(err))
		panic(err)
	}
	conn, err := Dial(addr)
	if err != nil {
		logger.LogService(serviceName).Error("dial error", zap.Error(err))
		panic(err)
	}
	return conn
}
