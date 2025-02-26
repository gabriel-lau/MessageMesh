package models

type FirstMessage struct {
	PeerID1      string `json:"peerID1"`
	PeerID2      string `json:"peerID2"`
	SymetricKey1 string `json:"symetricKey1"` // Symmetric key for peer 1 (Encrypted with peer 1's public key)
	SymetricKey2 string `json:"symetricKey2"` // Symmetric key for peer 2 (Encrypted with peer 2's public key)
}
