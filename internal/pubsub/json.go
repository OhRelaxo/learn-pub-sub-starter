package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	marsh, err := json.Marshal(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        marsh,
	})
	if err != nil {
		return err
	}
	return nil
}

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T)) error {
	queueChan, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}
	newChan, err := queueChan.Consume("", "", false, false, false, false, nil)
	if err != nil {
		return nil
	}

	go func() {
		for message := range newChan {
			var unmarMess T
			err := json.Unmarshal(message.Body, &unmarMess)
			if err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				return
			}
			handler(unmarMess)
			message.Ack(false)
		}
	}()

	return nil
}
