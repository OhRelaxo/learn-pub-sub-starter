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
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("successfully connected to RabbitMQ server")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("failed to get username: %v", err)
	}

	gameState := gamelogic.NewGameState(username)

	// subscribe to pause direct exchange
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		log.Fatalf("failed to subscribe to pause exchange: %v", err)
	}

	// subscirbe to topic exchange
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+gameState.GetUsername(), routing.ArmyMovesPrefix+".*", pubsub.Transient, handlerMove(gameState))
	if err != nil {
		log.Fatalf("failed to subscribe to topic exchange: %v", err)
	}

	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		command := userInput[0]
		switch command {
		case "spawn":
			err = gameState.CommandSpawn(userInput)
			if err != nil {
				log.Printf("failed to spawn unit: %v in location: %v; err: %v\n", userInput[1], userInput[2], err)
			}
		case "move":
			armyMove, err := gameState.CommandMove(userInput)
			if err != nil {
				log.Printf("failed to move army: %v\n", err)
				continue
			}
			err = pubsub.PublishJSON(publishCh, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+gameState.GetUsername(), armyMove)
			if err != nil {
				log.Printf("failed to move units: %v\n", err)
				continue
			}
			fmt.Printf("succesfully moved army to location: %v\n", armyMove.ToLocation)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Printf("error: command: %v was not found!\n", command)
			continue
		}
	}
}
