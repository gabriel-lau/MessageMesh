package models

type FirstMessage struct {
	PeerIDs      []string `json:"peerIDs"`
	SymetricKey0 string   `json:"symetricKey1"` // Symmetric key for peer[0] (Encrypted with peer[1]'s public key)
	SymetricKey1 string   `json:"symetricKey2"` // Symmetric key for peer[1] (Encrypted with peer[0]'s public key)
}

func (fm *FirstMessage) GetPeerIDs() []string {
	return fm.PeerIDs
}

func (fm *FirstMessage) GetSymetricKey(peerID string) string {
	peerIDs := fm.GetPeerIDs()
	if peerID == peerIDs[0] {
		return fm.SymetricKey0
	}
	return fm.SymetricKey1
}
