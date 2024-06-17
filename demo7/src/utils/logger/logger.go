package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"select-course/demo7/src/constant/config"
)

var (
	Logger *zap.Logger
)

func init() {
	// 初始化临时日志记录器以捕获早期错误
	tempLogger, _ := zap.NewDevelopment()
	defer tempLogger.Sync()

	// 控制台输出或生产环境配置
	if config.EnvCfg.ProjectMode == "prod" {
		// 确保日志目录存在
		if err := ensureLogDirectory(); err != nil {
			tempLogger.Error("创建日志目录失败", zap.Error(err))
			return
		}

		// 设置日志级别和选项
		level, options := configureLogLevelAndOptions()
		core := getZapCore(zapcore.NewJSONEncoder, level, options...)
		Logger = zap.New(core, options...)
	} else {
		Logger, _ = zap.NewDevelopment()
	}
}

// ensureLogDirectory 确保日志目录存在
func ensureLogDirectory() error {
	return os.MkdirAll(config.EnvCfg.LoggerDir, os.ModePerm)
}

// configureLogLevelAndOptions 根据配置设置日志级别和附加选项
func configureLogLevelAndOptions() (zapcore.Level, []zap.Option) {
	var level zapcore.Level
	var options []zap.Option

	switch config.EnvCfg.LoggerLevel {
	case "DEBUG":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(zap.DebugLevel))
	case "INFO":
		level = zap.InfoLevel
	case "WARN":
		level = zap.WarnLevel
	case "ERROR":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	default:
		level = zap.InfoLevel
	}
	// 添加调用者信息
	options = append(options, zap.AddCaller())
	return level, options
}

// getZapCore 根据配置创建zapcore.Core实例
func getZapCore(encoderFunc func(zapcore.EncoderConfig) zapcore.Encoder, level zapcore.Level, opts ...zap.Option) zapcore.Core {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	encoder := encoderFunc(encoderConfig)
	writer := getLogWriter()
	return zapcore.NewCore(encoder, writer, level)
}

// getLogWriter 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.EnvCfg.LoggerDir + "/" + config.EnvCfg.LoggerName,
		MaxSize:    config.EnvCfg.LoggerMaxSize,
		MaxBackups: config.EnvCfg.LoggerMaxBackups,
		MaxAge:     config.EnvCfg.LoggerMaxAge,
		Compress:   true,
	})
}

// LogService 返回一个日志记录器，带有服务器名称
func LogService(name string) *zap.Logger {
	return Logger.With(
		zap.String("server", name),
	)
}
