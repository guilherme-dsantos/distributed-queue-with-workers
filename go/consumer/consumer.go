// Package consumer consumer/consumer.go
package consumer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/guilherme-dsantos/distributed-queue-with-workers/go/message"
)

type Consumer struct {
	workerURL string
	handler   func(msg message.Message)
}

func New(workerURL string, handler func(msg message.Message)) *Consumer {
	return &Consumer{
		workerURL: workerURL,
		handler:   handler,
	}
}

func (c *Consumer) Start() {
	for {
		msg, err := c.receive()
		if err != nil {
			// If queue is empty, wait before trying again
			time.Sleep(5 * time.Second)
			continue
		}

		c.handler(msg)
	}
}

func (c *Consumer) receive() (message.Message, error) {
	resp, err := http.Post(c.workerURL+"/dequeue", "application/json", nil)
	if err != nil {
		return message.Message{}, fmt.Errorf("receive message: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusNoContent {
		return message.Message{}, fmt.Errorf("queue empty")
	}

	if resp.StatusCode != http.StatusOK {
		return message.Message{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var msg message.Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return message.Message{}, fmt.Errorf("decode message: %w", err)
	}

	return msg, nil
}
