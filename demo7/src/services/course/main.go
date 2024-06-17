package main

import (
	"context"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/oklog/run"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"select-course/demo7/src/constant/config"
	"select-course/demo7/src/constant/services"
	"select-course/demo7/src/rpc/course"
	"select-course/demo7/src/storage/database"
	"select-course/demo7/src/utils/breaker"
	"select-course/demo7/src/utils/consumer"
	"select-course/demo7/src/utils/discovery"
	"select-course/demo7/src/utils/limiter"
	"select-course/demo7/src/utils/local"
	"select-course/demo7/src/utils/logger"
	"select-course/demo7/src/utils/promet"
	"select-course/demo7/src/utils/tracing"
	"syscall"
)

func main() {
	// -------------------- init --------------------
	// tracing init
	tracer, closer := tracing.Init(services.CourseRpcServerName)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// metrics init
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
		grpcprom.WithServerCounterOptions(
			grpcprom.WithConstLabels(map[string]string{
				"service": services.CourseRpcServerName,
				"version": "v1",
				"env":     config.EnvCfg.ProjectMode,
			}),
		),
	)
	client := prometheus.NewPedanticRegistry()
	client.MustRegister(srvMetrics)

	// rpc init
	courseService := &Course{}
	rpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_opentracing.UnaryServerInterceptor(),
			srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(promet.ExtractContext)),
		),
		grpc.ChainStreamInterceptor(
			srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(promet.ExtractContext)),
		),
	)
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
	if err := local.InitLocal(); err != nil {
		logger.Logger.Error("local init error", zap.Error(err))
		panic(err)
	}

	// init mq
	if err := consumer.InitSelectListener(); err != nil {
		logger.Logger.Error("SelectConsumer init error for: %v", zap.Error(err))
		panic(err)

	}

	// init sentinel
	if err := sentinel.InitWithConfigFile("./sentinel.yml"); err != nil {
		logger.Logger.Error("sentinel init error for: %v", zap.Error(err))
		panic(err)
	}
	// load breaker
	circuitbreaker.RegisterStateChangeListeners(&breaker.StateChangeTestListener{}) //
	if _, err = circuitbreaker.LoadRules(breaker.ErrorCountRules); err != nil {
		logger.Logger.Error("breaker init error for: %v", zap.Error(err))
		panic(err)
	}
	// load limiter
	if _, err := hotspot.LoadRules(limiter.LimitRules); err != nil {
		logger.Logger.Error("limiter init error for: %v", zap.Error(err))
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
	httpSrv := &http.Server{Addr: fmt.Sprintf("%s:%d", config.EnvCfg.BaseHost, services.MetricsPort)}

	// collect metrics
	g.Add(func() error {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.HandlerFor(
			client,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		httpSrv.Handler = m
		return httpSrv.ListenAndServe()
	}, func(err error) {
		logger.LogService(services.CourseRpcServerName).Error("http server stopped", zap.Error(err))
		httpSrv.Close()
	})
	// back consumer
	g.Add(func() error {
		return consumer.SelectConsumer.Consumer()
	}, func(err error) {
		logger.LogService(services.CourseRpcServerName).Error("consumer stopped", zap.Error(err))
		consumer.SelectConsumer.Close()
	})
	// signal
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
	if err := g.Run(); err != nil {
		logger.LogService(services.CourseRpcServerName).Error("run error", zap.Error(err))
		if err := discovery.Consul.Deregister(context.Background(), services.CourseRpcServerName); err != nil {
			logger.LogService(services.CourseRpcServerName).Error("deregister error", zap.Error(err))
		}
		os.Exit(1)
	}
}
