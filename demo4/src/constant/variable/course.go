package variable

import "github.com/streadway/amqp"

const (
	SelectExchange   = "select-course:select_exchange"
	SelectRoutingKey = "select-course:select_routing_key"
	SelectKind       = amqp.ExchangeDirect
	SelectQueue      = "select-course:select_queue"

	// 声明死信队列

	DeadExchange   = "select-course:dead_exchange"
	DeadRoutingKey = "select-course:dead_routing_key"
	DeadKind       = amqp.ExchangeDirect
	DeadQueue      = "select-course:dead_queue"
)
