// producer example
package main

import (
	"log"
	"os"

	"github.com/guilherme-dsantos/distributed-queue-with-workers/go/producer"
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
	p := producer.New(goDotEnvVariable("WORKER_URL"))

	id, err := p.Send(map[string]string{
		"task": "send_email",
		"to":   "user@example.com",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sent message with ID: %s", id)
}
