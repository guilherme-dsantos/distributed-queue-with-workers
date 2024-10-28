package message

// Message represents a queue message
type Message struct {
    ID      string      `json:"id,omitempty"`
    Payload interface{} `json:"payload"`
}
