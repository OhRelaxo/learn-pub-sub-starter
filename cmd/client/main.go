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
	fmt.Println("Starting Peril client...")

	conn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("failed to connect to amqp server: %v", err)
	}
	defer conn.Close()
	log.Println("successfully connected to server")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("failed to get username: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Printf("could not subscribe to pause: %v\n", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gameState := gamelogic.NewGameState(username)

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
