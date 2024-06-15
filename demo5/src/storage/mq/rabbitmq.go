package mq

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/streadway/amqp"
	"select-course/demo5/src/constant/config"
	"select-course/demo5/src/utils/logger"
	"time"
)

var Client *amqp.Connection

func init() {
	dns := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		config.EnvCfg.RabbitMQUser,
		config.EnvCfg.RabbitMQPassword,
		config.EnvCfg.RabbitMQHost,
		config.EnvCfg.RabbitMQPort,
		config.EnvCfg.RabbitMQVhost,
	)
	// 等待rabbitmq启动完成
	if err := retry.Do(func() error {
		conn, err := amqp.Dial(dns)
		if err != nil {
			return err
		}
		Client = conn
		return nil
	}, retry.Attempts(5), retry.Delay(time.Second)); err != nil {
		logger.Logger.Error("rabbitmq init fail")
		panic(err)
	}
	logger.Logger.Info("rabbitmq init success")

}
