package config

import (
	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

var EnvCfg envConfig

type envConfig struct {
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
}

func init() {
	if err := godotenv.Load("demo1/.env.dev"); err != nil {
		log.Fatalln("read .env file failed")
	}
	EnvCfg = envConfig{}

	if err := env.Parse(&EnvCfg); err != nil {
		panic("Can not parse env from file system, please check the env.")
	}
}