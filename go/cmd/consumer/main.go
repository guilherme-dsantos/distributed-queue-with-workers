// consumer example
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/guilherme-dsantos/distributed-queue-with-workers/go/consumer"
	"github.com/guilherme-dsantos/distributed-queue-with-workers/go/message"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load("../worker.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	// Define message handler
	handleMessage := func(msg message.Message) {
		fmt.Printf("Received message: %+v\n", msg)
		// Process the message...
	}

	// Create and start consumer
	c := consumer.New(goDotEnvVariable("WORKER_URL"), handleMessage)
	c.Start()
}
