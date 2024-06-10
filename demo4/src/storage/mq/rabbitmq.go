package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"select-course/demo4/src/constant/config"
	"select-course/demo4/src/utils/logger"
)

var Client *amqp.Connection

func init() {
	// "amqp://" + rabbit.Username + ":" + rabbit.Password + "@" + rabbit.Host + ":" + strconv.Itoa(rabbit.Port) + "/" + rabbit.Vhost
	dns := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		config.EnvCfg.RabbitMQUser,
		config.EnvCfg.RabbitMQPassword,
		config.EnvCfg.RabbitMQHost,
		config.EnvCfg.RabbitMQPort,
		config.EnvCfg.RabbitMQVhost,
	)
	conn, err := amqp.Dial(dns)
	if err != nil {
		panic(err)
	}

	Client = conn

	logger.Logger.Info("rabbitmq init success")

}
