// producer/producer.go
package producer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/guilherme-dsantos/distributed-queue-with-workers/go/message"
)

type Producer struct {
	workerURL string
}

func New(workerURL string) *Producer {
	return &Producer{workerURL: workerURL}
}

func (p *Producer) Send(payload interface{}) (string, error) {
	msg := message.Message{
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshal message: %w", err)
	}

	resp, err := http.Post(p.workerURL+"/enqueue", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.ID, nil
}
