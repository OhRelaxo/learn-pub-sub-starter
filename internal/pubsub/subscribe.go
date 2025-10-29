package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype = int

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func subscribe[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype, unmarshaller func([]byte) (T, error)) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}
	messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	go func() {
		defer ch.Close()
		for message := range messages {
			decodedMessage, err := unmarshaller(message.Body)
			if err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				return
			}
			acktype := handler(decodedMessage)
			switch acktype {
			case Ack:
				message.Ack(false)
			case NackRequeue:
				message.Nack(false, true)
			case NackDiscard:
				message.Nack(false, false)
			}
		}
	}()
	return nil
}

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {
	return subscribe(conn, exchange, queueName, key, queueType, handler, func(b []byte) (T, error) {
		var decodedMessage T
		err := json.Unmarshal(b, &decodedMessage)
		if err != nil {
			return decodedMessage, err
		}
		return decodedMessage, nil
	})
}

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {
	return subscribe(conn, exchange, queueName, key, queueType, handler, func(b []byte) (T, error) {
		buffer := bytes.NewBuffer(b)
		decoder := gob.NewDecoder(buffer)
		var decodedMessage T
		err := decoder.Decode(&decodedMessage)
		if err != nil {
			return decodedMessage, err
		}
		return decodedMessage, nil
	})
}
