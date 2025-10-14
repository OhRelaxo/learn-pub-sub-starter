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

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	connChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	var isDurable bool
	var autoDelete bool
	var exclusive bool
	switch queueType {
	case 0:
		isDurable = true
		autoDelete = false
		exclusive = false
	case 1:
		isDurable = false
		autoDelete = true
		exclusive = true
	default:
		return nil, amqp.Queue{}, fmt.Errorf("please use a valid value for queueType")
	}
	queue, err := connChan.QueueDeclare(queueName, isDurable, autoDelete, exclusive, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	err = connChan.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return connChan, queue, nil
}
