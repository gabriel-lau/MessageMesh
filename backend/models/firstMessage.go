package models

type FirstMessage struct {
	PeerIDs      []string `json:"peerIDs"`
	SymetricKey0 []byte   `json:"symetricKey1"` // Symmetric key for peer[0] (Encrypted with peer[1]'s public key)
	SymetricKey1 []byte   `json:"symetricKey2"` // Symmetric key for peer[1] (Encrypted with peer[0]'s public key)
}

func (fm *FirstMessage) GetSymetricKey(peerID string) []byte {
	if peerID == fm.PeerIDs[0] {
		return fm.SymetricKey0
	}
	return fm.SymetricKey1
}
