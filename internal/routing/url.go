package routing

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetServerUrl() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("RABBITMQ_DEFAULT_USER")
	password := os.Getenv("RABBITMQ_DEFAULT_PASS")
	return fmt.Sprintf("amqp://%v:%v@localhost:5672/", username, password)
}
