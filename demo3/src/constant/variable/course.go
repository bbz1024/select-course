package variable

import "github.com/streadway/amqp"

const (
	SelectExchange   = "select-course:select_exchange"
	SelectRoutingKey = "select-course:select_routing_key"
	SelectKind       = amqp.ExchangeDirect
	SelectQueue      = "select-course:select_queue"
)
