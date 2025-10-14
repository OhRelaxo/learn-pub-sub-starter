package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const serverUrl = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("failed to connect to amqp server: %v", err)
	}
	defer conn.Close()
	log.Println("successfully connected to server")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("failed to create user: %v", err)
	}
	_, _, err = pubsub.DeclareAndBind(conn, "peril_direct", fmt.Sprintf("pause.%v", username), "pause", pubsub.Transient)
	if err != nil {
		log.Println(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Println("shutting down")
}
