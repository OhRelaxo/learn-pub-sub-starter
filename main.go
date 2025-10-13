package main

import (
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const serverUrl = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("faild to connect to amqp server: %v", err)
	}
	defer conn.Close()
	log.Println("successfully connected to server")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Println("shutting down")
}
