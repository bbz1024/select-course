package database

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"select-course/demo1/src/constant/config"
	"time"
)

var Client *sql.DB

func init() {
	fmt.Println(config.EnvCfg.MySqlHOST)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.EnvCfg.MysqlUser,
			config.EnvCfg.MysqlPassword,
			config.EnvCfg.MySqlHOST,
			config.EnvCfg.MysqlPort,
			config.EnvCfg.MysqlDatabase,
		),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置

	}), &gorm.Config{
		Logger: getGormLogger(), // 打印日志
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表明不加s
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	Client, err = db.DB()

	Client.SetMaxIdleConns(config.EnvCfg.MysqlMaxIdleConns) // 设置连接池，空闲
	Client.SetMaxOpenConns(config.EnvCfg.MysqlMaxOpenConns) // 打开
	Client.SetConnMaxLifetime(time.Second * 30)

}
func getGormLogger() logger.Interface {
	var logMode logger.LogLevel
	switch config.EnvCfg.MysqlLogLevel {
	case "SILENT":
		logMode = logger.Silent
	case "ERROR":
		logMode = logger.Error
	case "WARN":
		logMode = logger.Warn
	case "INFO":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}
	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
		LogLevel:                  logMode,                // 日志级别
		IgnoreRecordNotFoundError: false,                  // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  false,                  // 禁用彩色打印
	})
}

// 自定义 gorm Writer
func getGormLogWriter() logger.Writer {
	var writer io.Writer
	writer = os.Stdout
	// 是否启用日志文件
	/*
		if global.Config.Mysql.EnableFileLogWriter {
			// 自定义 Writer
			writer = &lumberjack.Logger{
				Filename:   global.Config.Log.RootDir + "/" + global.Config.Mysql.LogFilename,
				MaxSize:    global.Config.Log.MaxSize,
				MaxBackups: global.Config.Log.MaxBackups,
				MaxAge:     global.Config.Log.MaxAge,
				Compress:   global.Config.Log.Compress,
			}
		}

	*/
	return log.New(writer, "\r\n", log.LstdFlags)
}
