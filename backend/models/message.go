package models

type Message struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
