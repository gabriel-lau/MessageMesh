package models

type FirstMessage struct {
	PeerID1      string `json:"peerID1"`
	PeerID2      string `json:"peerID2"`
	SymetricKey1 string `json:"symetricKey1"`
	SymetricKey2 string `json:"symetricKey2"`
}
