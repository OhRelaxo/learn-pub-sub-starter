package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	serverUrl := routing.GetServerUrl()
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("faild to connect to amqp server: %v", err)
	}
	defer conn.Close()
	log.Println("successfully connected to server")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to create a channel from estableshed connection: %v", err)
	}

	_, topicQueue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable)
	if err != nil {
		log.Fatalf("Could not open topic queue: %v", err)
	}
	log.Printf("Queue %v declared and bound!\n", topicQueue.Name)

	gamelogic.PrintServerHelp()

	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		command := userInput[0]
		switch command {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Printf("failed to pause game: %v\n", err)
			}
		case "resume":
			fmt.Println("Publishing resume game state")
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Printf("failed to pause game: %v\n", err)
			}
		case "quit":
			log.Println("exiting...")
			return
		default:
			log.Printf("command: %v not found!", command)
		}
	}
}
