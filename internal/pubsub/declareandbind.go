package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var isDurable bool
	var autoDelete bool
	var exclusive bool

	switch queueType {
	case Durable:
		isDurable = true
		autoDelete = false
		exclusive = false
	case Transient:
		isDurable = false
		autoDelete = true
		exclusive = true
	default:
		return nil, amqp.Queue{}, fmt.Errorf("please use a valid value for queueType")
	}

	table := amqp.Table{
		"x-dead-letter-exchange": "peril_dlx",
	}

	queue, err := ch.QueueDeclare(queueName, isDurable, autoDelete, exclusive, false, table)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}
