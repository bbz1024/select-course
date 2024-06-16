package config

import (
	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var EnvCfg envConfig

type envConfig struct {
	BaseHost string `env:"BASE_HOST" envDefault:"localhost"`
	// project
	ProjectName string `env:"PROJECT_NAME" envDefault:"SelectCourse"`
	ServerPort  int    `env:"PROJECT_PORT" envDefault:"8080"`
	ProjectMode string `env:"PROJECT_MODE" envDefault:"dev"`
	// logger
	LoggerLevel string `env:"PROJECT_LOG_LEVEL" envDefault:"DEBUG"`

	LoggerDir        string `env:"PROJECT_LOG_DIR" envDefault:"./logs"`
	LoggerName       string `env:"PROJECT_LOG_NAME" envDefault:"app.log"`
	LoggerMaxBackups int    `env:"PROJECT_LOG_MAX_BACKUPS" envDefault:"3"`
	LoggerMaxSize    int    `env:"PROJECT_LOG_MAX_SIZE" envDefault:"500"`
	LoggerMaxAge     int    `env:"PROJECT_LOG_MAX_AGE" envDefault:"28"`

	// mysql
	MySqlHOST         string `env:"MYSQL_ROOT_HOST" envDefault:"localhost"`
	MysqlPort         int    `env:"MYSQL_ROOT_PORT" envDefault:"3306"`
	MysqlUser         string `env:"MYSQL_ROOT_USER" envDefault:"root"`
	MysqlPassword     string `env:"MYSQL_ROOT_PASSWORD" envDefault:"root"`
	MysqlDatabase     string `env:"MysqlDatabase" envDefault:"test"`
	MysqlLogLevel     string `env:"MysqlLogLevel" envDefault:"debug"`
	MysqlMaxIdleConns int    `env:"MysqlMaxIdleConns" envDefault:"10"`
	MysqlMaxOpenConns int    `env:"MysqlMaxOpenConns" envDefault:"100"`

	// redis
	RedisHost           string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort           int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisPwd            string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDb             int    `env:"RedisDb" envDefault:"0"`
	RedisMaxIdleConns   int    `env:"RedisMaxIdleConns" envDefault:"10"`
	RedisMaxActiveConns int    `env:"RedisMaxActiveConns" envDefault:"100"`

	// mq
	RabbitMQHost     string `env:"RABBITMQ_DEFAULT_HOST" envDefault:"localhost"`
	RabbitMQPort     int    `env:"RABBITMQ_DEFAULT_PORT" envDefault:"5672"`
	RabbitMQUser     string `env:"RABBITMQ_DEFAULT_USER" envDefault:"guest"`
	RabbitMQPassword string `env:"RABBITMQ_DEFAULT_PASS" envDefault:"guest"`
	RabbitMQVhost    string `env:"RABBITMQ_DEFAULT_VHOST" envDefault:"/"`

	// jaeger
	JaegerHost string `env:"JAEGER_AGENT_HOST" envDefault:"localhost"`
	JaegerPort int    `env:"JAEGER_AGENT_PORT" envDefault:"6831"`

	// consul
	ConsulHost string `env:"CONSUL_HOST" envDefault:"localhost"`
	ConsulPort string `env:"CONSUL_PORT" envDefault:"8500"`
}

var path = ".env.dev"

//var path = ".env.prod"

func init() {
	prod := os.Getenv("PROJECT_MODE")
	if prod == "prod" {
		path = ".env"
	}
	if err := godotenv.Load(path); err != nil {
		log.Fatalln("read .env file failed")
	}
	EnvCfg = envConfig{}
	if err := env.Parse(&EnvCfg); err != nil {
		panic("Can not parse env from file system, please check the env.")
	}

}
