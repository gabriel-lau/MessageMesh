package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"fmt"
	"sort"
	"time"
)

func (network *Network) ConnectToNetwork() {
	debug.Log("server", "This may take upto 30 seconds.")

	// Create a new P2PHost
	network.P2pService = NewP2PService()
	debug.Log("server", "Completed P2P Setup")

	// Connect to peers with the chosen discovery method
	network.P2pService.AdvertiseConnect()
	// network.P2p.AnnounceConnect()
	network.Progress.NetworkConnected <- true
	debug.Log("server", "Connected to Service Peers")
	// Join the chat room
	network.PubSubService, _ = JoinPubSub(network.P2pService)
	network.Progress.PubSubJoined <- true
	debug.Log("server", "Joined the PubSub")
	// Wait for network setup to complete
	time.Sleep(time.Second * 5)
	debug.Log("server", "Connected to Service Peers")

	// Print my peer ID
	debug.Log("server", fmt.Sprintf("My Peer ID: %s", network.PubSubService.SelfID()))

	// Print my multiaddress
	debug.Log("server", fmt.Sprintf("My Multiaddress: %s", network.P2pService.AllNodeAddr()))

	network.ConsensusService, _ = StartConsensus(network)
	network.Progress.BlockchainLoaded <- true
	debug.Log("server", "Blockchain loaded")
}

func (network *Network) SendMessage(message string, receiver string) {
	sender := network.PubSubService.SelfID().String() // Self ID
	network.PubSubService.Outbound <- models.Message{
		Sender:    sender,
		Receiver:  receiver,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func (network *Network) SendEncryptedMessage(message string, receiver string) {
	sender := network.PubSubService.SelfID().String() // Self ID
	peerIDs := []string{sender, receiver}
	sort.Strings(peerIDs)

	// Encrypt the message with the symmetric key
	encryptedMessage, err := network.EncryptMessage(message, receiver)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error encrypting message for %s: %s", receiver, err.Error()))
		return
	}
	network.SendMessage(encryptedMessage, receiver)
}

func (network *Network) EncryptMessage(message string, receiver string) (string, error) {
	sender := network.PubSubService.SelfID().String() // Self ID
	peerIDs := []string{sender, receiver}
	sort.Strings(peerIDs)

	// Get the symmetric key for the two peers if it exists in the database
	symmetricKey, err := GetSymmetricKey(peerIDs)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error getting symmetric key for %s and %s: %s", peerIDs[0], peerIDs[1], err.Error()))
		return "", err
	}
	// If the symmetric key is not found check the blockchain for a first message
	if symmetricKey == nil {
		debug.Log("server", fmt.Sprintf("Symmetric key not found for %s and %s", peerIDs[0], peerIDs[1]))
		// Check if the firstMessage is shared between the two peers in the blockchain
		firstMessage := network.ConsensusService.Blockchain.CheckPeerFirstMessage(peerIDs)
		if firstMessage != nil {
			debug.Log("server", fmt.Sprintf("First message found for %s and %s", peerIDs[0], peerIDs[1]))
			symmetricKey = []byte(firstMessage.GetSymetricKey(sender))
		}
	}

	// Encrypt the message with the symmetric key
	encryptedMessage, err := EncryptWithSymmetricKey([]byte(message), symmetricKey)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error encrypting message for %s: %s", receiver, err.Error()))
		return "", err
	}
	debug.Log("server", fmt.Sprintf("Encrypted message for %s", receiver))
	return string(encryptedMessage), nil
}

func (network *Network) SendFirstMessage(peerIDs []string) (models.FirstMessage, error) {
	// Generate a symmetric key
	symmetricKey, err := GenerateSymmetricKey(32)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error generating symmetric key: %s", err.Error()))
		return models.FirstMessage{}, err
	}

	// Save the symmetric key to the database
	err = SaveSymmetricKey(symmetricKey, peerIDs)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error saving symmetric key: %s", err.Error()))
		return models.FirstMessage{}, err
	}

	peerID0 := peerIDs[0]
	peerID1 := peerIDs[1]

	// Encrypt the symmetric key with the peer[0]'s public key
	encryptedSymmetricKey0, err := EncryptForPeer(network.P2pService, symmetricKey, peerID0)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error encrypting symmetric key for %s: %s", peerID0, err.Error()))
		return models.FirstMessage{}, err
	}
	debug.Log("server", fmt.Sprintf("Encrypted symmetric key for %s", peerID0))

	// Encrypt the symmetric key with the peer[1]'s public key
	encryptedSymmetricKey1, err := EncryptForPeer(network.P2pService, symmetricKey, peerID1)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error encrypting symmetric key for %s: %s", peerID1, err.Error()))
		return models.FirstMessage{}, err
	}
	debug.Log("server", fmt.Sprintf("Encrypted symmetric key for %s", peerID1))

	// Create the first message
	firstMessage := models.FirstMessage{
		PeerIDs:      peerIDs,
		SymetricKey0: string(encryptedSymmetricKey0),
		SymetricKey1: string(encryptedSymmetricKey1),
	}
	network.PubSubService.Outbound <- firstMessage
	debug.Log("server", fmt.Sprintf("First message sent to %s and %s", peerID0, peerID1))
	return firstMessage, nil
}
