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
	const serverUrl = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("faild to connect to amqp server: %v", err)
	}
	defer conn.Close()
	log.Println("successfully connected to server")

	gamelogic.PrintServerHelp()

	newChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to create a channel from estableshed connection: %v", err)
	}

	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		command := userInput[0]
		switch command {
		case "pause":
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Printf("failed to pause game: %v", err)
			}
		case "resume":
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Printf("failed to pause game: %v", err)
			}
		case "quit":
			log.Println("exiting...")
			return
		default:
			log.Printf("command: %v not found!", command)
		}
	}
}
