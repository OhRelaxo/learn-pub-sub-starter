package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/rabbitmq/amqp091-go"
)

func handleWar(gs *gamelogic.GameState, ch *amqp091.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(players gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(players)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			err := pubsub.PublishGameLog(ch, gs.GetUsername(), fmt.Sprintf("%v won a war against %v", winner, loser))
			if err != nil {
				log.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			err := pubsub.PublishGameLog(ch, gs.GetUsername(), fmt.Sprintf("%v won a war against %v", winner, loser))
			if err != nil {
				log.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			err := pubsub.PublishGameLog(ch, gs.GetUsername(), fmt.Sprintf("A war between %v and %v resulted in a draw", winner, loser))
			if err != nil {
				log.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Printf("error got unexpected war outcome: %v", outcome)
			return pubsub.NackDiscard
		}
	}
}
