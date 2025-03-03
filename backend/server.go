package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
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
	msg := models.Message{
		Sender:    sender,
		Receiver:  receiver,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Marshal message to JSON
	messageJSON, err := json.Marshal(msg)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error marshaling message: %s", err.Error()))
		return
	}
	debug.Log("server", fmt.Sprintf("Message marshaled: %s", string(messageJSON)))

	network.PubSubService.Outbound <- MessageEnvelope{
		Type: "Message",
		Data: messageJSON,
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
	debug.Log("server", fmt.Sprintf("Sending encrypted message: %s", encryptedMessage))
	network.SendMessage(encryptedMessage, receiver)
}

func (network *Network) EncryptMessage(message string, receiver string) (string, error) {
	sender := network.PubSubService.SelfID().String() // Self ID
	peerIDs := []string{sender, receiver}
	sort.Strings(peerIDs)

	// Get the symmetric key for the two peers if it is saved
	symmetricKey, err := GetSymmetricKey(peerIDs)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error getting symmetric key for %s and %s: %s", peerIDs[0], peerIDs[1], err.Error()))
	}
	// If the symmetric key is not found check the blockchain for a first message
	if symmetricKey == nil {
		debug.Log("server", fmt.Sprintf("Symmetric key not found for %s and %s", peerIDs[0], peerIDs[1]))
		// Check if the firstMessage is shared between the two peers in the blockchain
		firstMessage := network.ConsensusService.Blockchain.CheckPeerFirstMessage(peerIDs)
		keyPair, err := ReadKeyPair()
		if err != nil {
			debug.Log("server", fmt.Sprintf("Error reading key pair: %s", err.Error()))
			return "", err
		}
		// If the first message is found, decrypt the symmetric key with the private key
		if firstMessage != nil {
			debug.Log("server", fmt.Sprintf("First message found for %s and %s", peerIDs[0], peerIDs[1]))
			symmetricKey, err = keyPair.DecryptWithPrivateKey(firstMessage.GetSymetricKey(sender))
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error decrypting symmetric key for %s: %s", sender, err.Error()))
				return "", err
			}
		} else {
			// If the first message is not found, send a first message and decrypt the symmetric key with the private key
			debug.Log("server", fmt.Sprintf("First message not found for %s and %s", peerIDs[0], peerIDs[1]))
			firstMessage, err := network.SendFirstMessage(peerIDs, receiver)
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error sending first message to %s and %s: %s", peerIDs[0], peerIDs[1], err.Error()))
				return "", err
			}
			symmetricKey, err = keyPair.DecryptWithPrivateKey(firstMessage.GetSymetricKey(sender))
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error decrypting symmetric key for %s: %s", sender, err.Error()))
				return "", err
			}
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

// Decrypt message with the symmetric key
func (network *Network) DecryptMessage(message string, sender string) (string, error) {
	receiver := network.PubSubService.SelfID().String()
	peerIDs := []string{sender, receiver}
	sort.Strings(peerIDs)
	// Get the symmetric key for the two peers if it is saved
	symmetricKey, err := GetSymmetricKey(peerIDs)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error getting symmetric key for %s and %s: %s", peerIDs[0], peerIDs[1], err.Error()))
	}
	// If the symmetric key is not found check the blockchain for a first message
	if symmetricKey == nil {
		debug.Log("server", fmt.Sprintf("Symmetric key not found for %s and %s", peerIDs[0], peerIDs[1]))
		// Check if the firstMessage is shared between the two peers in the blockchain
		firstMessage := network.ConsensusService.Blockchain.CheckPeerFirstMessage(peerIDs)
		err = SaveSymmetricKey(symmetricKey, peerIDs)
		if err != nil {
			debug.Log("server", fmt.Sprintf("Error saving symmetric key: %s", err.Error()))
			return "", err
		}
		keyPair, err := ReadKeyPair()
		if err != nil {
			debug.Log("server", fmt.Sprintf("Error reading key pair: %s", err.Error()))
			return "", err
		}
		// If the first message is found, decrypt the symmetric key with the private key
		if firstMessage != nil {
			symmetricKey, err = keyPair.DecryptWithPrivateKey(firstMessage.GetSymetricKey(receiver))
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error getting symmetric key for %s and %s: %s", peerIDs[0], peerIDs[1], err.Error()))
			}
		} else {
			debug.Log("server", fmt.Sprintf("First message not found for %s and %s", peerIDs[0], peerIDs[1]))
			return "", fmt.Errorf("first message not found for %s and %s", peerIDs[0], peerIDs[1])
		}
	}

	// Decrypt the message with the symmetric key
	decryptedMessage, err := DecryptWithSymmetricKey([]byte(message), symmetricKey)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error decrypting message: %s", err.Error()))
		return "", err
	}
	return string(decryptedMessage), nil
}

func (network *Network) SendFirstMessage(peerIDs []string, receiver string) (models.FirstMessage, error) {
	receiverPeerID := peer.ID(receiver)
	// Check if the user is online
	if receiverPeerID == "" {
		debug.Log("server", fmt.Sprintf("User %s is not online", receiver))
		return models.FirstMessage{}, fmt.Errorf("user %s is not online", receiver)
	}
	// Generate a symmetric key
	symmetricKey, err := GenerateSymmetricKey(32)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error generating symmetric key: %s", err.Error()))
		return models.FirstMessage{}, err
	}
	debug.Log("server", fmt.Sprintf("Symmetric key generated for %s and %s", peerIDs[0], peerIDs[1]))

	// Save the symmetric key to the database
	err = SaveSymmetricKey(symmetricKey, peerIDs)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error saving symmetric key: %s", err.Error()))
		return models.FirstMessage{}, err
	}
	debug.Log("server", fmt.Sprintf("Symmetric key saved for %s and %s", peerIDs[0], peerIDs[1]))

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
		SymetricKey0: encryptedSymmetricKey0,
		SymetricKey1: encryptedSymmetricKey1,
	}

	// Marshal firstMessage to JSON first
	firstMessageJSON, err := json.Marshal(firstMessage)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error marshaling first message: %s", err.Error()))
		return models.FirstMessage{}, err
	}

	network.PubSubService.Outbound <- MessageEnvelope{
		Type: "FirstMessage",
		Data: firstMessageJSON,
	}
	debug.Log("server", fmt.Sprintf("First message sent to %s and %s", peerID0, peerID1))
	return firstMessage, nil
}
