package config

import (
	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

var EnvCfg envConfig

type envConfig struct {
	// project
	ServerPort  string `env:"ServerPort" envDefault:"8080"`
	ProjectMode string `env:"ProjectMode" envDefault:"dev"`
	// logger
	LoggerLevel string `env:"LoggerLevel" envDefault:"DEBUG"`

	// mysql
	MySqlHOST         string `env:"MysqlHOST" envDefault:"localhost"`
	MysqlPort         string `env:"MysqlPort" envDefault:"3306"`
	MysqlUser         string `env:"MysqlUser" envDefault:"root"`
	MysqlPassword     string `env:"MysqlPassword" envDefault:"root"`
	MysqlDatabase     string `env:"MysqlDatabase" envDefault:"test"`
	MysqlLogLevel     string `env:"MysqlLogLevel" envDefault:"debug"`
	MysqlMaxIdleConns int    `env:"MysqlMaxIdleConns" envDefault:"10"`
	MysqlMaxOpenConns int    `env:"MysqlMaxOpenConns" envDefault:"100"`

	// redis
	RedisHost           string `env:"RedisHost" envDefault:"localhost"`
	RedisPort           string `env:"RedisPort" envDefault:"6379"`
	RedisPwd            string `env:"RedisPwd" envDefault:""`
	RedisDb             int    `env:"RedisDb" envDefault:"0"`
	RedisMaxIdleConns   int    `env:"RedisMaxIdleConns" envDefault:"10"`
	RedisMaxActiveConns int    `env:"RedisMaxActiveConns" envDefault:"100"`
	// mqm

	RabbitMQHost     string `env:"RabbitMQHost" envDefault:"localhost"`
	RabbitMQPort     string `env:"RabbitMQPort" envDefault:"5672"`
	RabbitMQUser     string `env:"RabbitMQUser" envDefault:"guest"`
	RabbitMQPassword string `env:"RabbitMQPassword" envDefault:"guest"`
	RabbitMQVhost    string `env:"RabbitMQVhost" envDefault:"test"`
}

//var path = "../../.env.dev"

var path = "demo3/.env.dev"

func init() {

	if err := godotenv.Load(path); err != nil {
		log.Fatalln("read .env file failed")
	}
	EnvCfg = envConfig{}

	if err := env.Parse(&EnvCfg); err != nil {
		panic("Can not parse env from file system, please check the env.")
	}
}
