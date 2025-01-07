package models

import "github.com/libp2p/go-libp2p-core/peer"

type Message struct {
	Sender       peer.ID       `json:"sender"`
	Receiver     peer.ID       `json:"receiver"`
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
