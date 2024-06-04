package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"select-course/demo3/src/constant/config"
	"select-course/demo3/src/utils/logger"
)

var Client *amqp.Channel

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
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	Client = channel

	logger.Logger.Info("rabbitmq init success")

}
