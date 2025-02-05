package models

type FirstMessage struct {
	Sender      string `json:"sender"`
	Receiver    string `json:"receiver"`
	SymetricKey string `json:"symetricKey"`
}
