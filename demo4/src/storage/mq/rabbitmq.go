package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"select-course/demo4/src/constant/config"
	"select-course/demo4/src/utils/logger"
	"time"
)

var Client *amqp.Connection

func init() {
	// "amqp://" + rabbit.Username + ":" + rabbit.Password + "@" + rabbit.Host + ":" + strconv.Itoa(rabbit.Port) + "/" + rabbit.Vhost
	dns := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		config.EnvCfg.RabbitMQUser,
		config.EnvCfg.RabbitMQPassword,
		config.EnvCfg.RabbitMQHost,
		config.EnvCfg.RabbitMQPort,
		config.EnvCfg.RabbitMQVhost,
	)
	// 等待rabbitmq启动完成
	for {
		conn, err := amqp.Dial(dns)
		if err != nil {
			time.Sleep(time.Millisecond * 100)
			logger.Logger.Info("rabbitmq not ready, retry...", err)
			continue
		}
		err = conn.Close()
		if err != nil {
			panic(err)
		}
		logger.Logger.Info("rabbitmq ready")
		break
	}
	conn, err := amqp.Dial(dns)

	if err != nil {
		panic(err)
	}

	Client = conn

	logger.Logger.Info("rabbitmq init success")

}
