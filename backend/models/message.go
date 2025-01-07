package models

type Message struct {
	Sender       string        `json:"sender"`
	Receiver     string        `json:"receiver"`
	Message      string        `json:"message"`
	Timestamp    string        `json:"timestamp"`
	FirstMessage *FirstMessage `json:"firstMessage"`
}

type FirstMessage struct {
	SymetricKey string `json:"symetricKey"`
}

func (m Message) GetType() string {
	return "Message"
}
