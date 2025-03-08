package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"MessageMesh/monitoring"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

func (network *Network) ConnectToNetwork() {
	debug.Log("server", "This may take upto 30 seconds.")

	// Initialize system monitor
	monitor, err := monitoring.NewSystemMonitor()
	if err != nil {
		debug.Log("server", fmt.Sprintf("Failed to initialize system monitor: %s", err.Error()))
	} else {
		// Start periodic monitoring
		go network.runMonitoring(monitor)
	}

	// Create a new P2PHost
	network.P2pService = NewP2PService()
	debug.Log("server", "Completed P2P Setup")

	// Connect to peers with the chosen discovery method
	network.P2pService.AdvertiseConnect()
	// network.P2p.AnnounceConnect()
	debug.Log("server", "Connected to Service Peers")
	// Join the chat room
	network.PubSubService, _ = JoinPubSub(network.P2pService)
	debug.Log("server", "Joined the PubSub")
	// Wait for network setup to complete
	time.Sleep(time.Second * 5)
	debug.Log("server", "Connected to Service Peers")

	debug.Log("server", fmt.Sprintf("My Peer ID: %s", network.PubSubService.SelfID()))
	debug.Log("server", fmt.Sprintf("My Multiaddress: %s", network.P2pService.AllNodeAddr()))

	network.ConsensusService, _ = StartConsensus(network)
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
		debug.Log("server", fmt.Sprintf("Symmetric key not found in keyMap for %s and %s", peerIDs[0], peerIDs[1]))
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
	// Convert to base64 string
	base64Message := base64.StdEncoding.EncodeToString(encryptedMessage)
	debug.Log("server", fmt.Sprintf("Encrypted message for %s", receiver))
	return base64Message, nil
}

// Decrypt message with the symmetric key
func (network *Network) DecryptMessage(message string, peerIDs []string) (string, error) {
	sort.Strings(peerIDs)
	// Get the symmetric key for the two peers if it is saved
	symmetricKey, err := GetSymmetricKey(peerIDs)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error getting symmetric key for %s and %s: %s", peerIDs[0], peerIDs[1], err.Error()))
	}
	// If the symmetric key is not found check the blockchain for a first message
	if symmetricKey == nil {
		debug.Log("server", fmt.Sprintf("Symmetric key not found in keyMap for %s and %s", peerIDs[0], peerIDs[1]))
		// Check if the firstMessage is shared between the two peers in the blockchain
		firstMessage := network.ConsensusService.Blockchain.CheckPeerFirstMessage(peerIDs)
		if firstMessage == nil {
			debug.Log("server", fmt.Sprintf("First message not found for %s and %s", peerIDs[0], peerIDs[1]))
			return "", fmt.Errorf("first message not found for %s and %s", peerIDs[0], peerIDs[1])
		}
		keyPair, err := ReadKeyPair()
		if err != nil {
			debug.Log("server", fmt.Sprintf("Error reading key pair: %s", err.Error()))
			return "", err
		}
		debug.Log("server", "Reading key pair")
		// If the first message is found, decrypt the symmetric key with the private key
		if firstMessage != nil {
			symmetricKey, err = keyPair.DecryptWithPrivateKey(firstMessage.GetSymetricKey(network.PubSubService.selfid.String()))
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error decrypting symmetric key for %s and %s: %s", peerIDs[0], peerIDs[1], err.Error()))
			}
			debug.Log("server", fmt.Sprintf("Decrypted symmetric key for %s and %s", peerIDs[0], peerIDs[1]))
			err = SaveSymmetricKey(symmetricKey, peerIDs)
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error saving symmetric key: %s", err.Error()))
				return "", err
			}
			debug.Log("server", fmt.Sprintf("Saved symmetric key for %s and %s", peerIDs[0], peerIDs[1]))
		} else {
			debug.Log("server", fmt.Sprintf("First message not found for %s and %s", peerIDs[0], peerIDs[1]))
			return "", fmt.Errorf("first message not found for %s and %s", peerIDs[0], peerIDs[1])
		}
	}

	// Decode from base64 first
	encryptedBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error decoding base64 message: %s", err.Error()))
		return "", err
	}

	// Decrypt the message with the symmetric key
	decryptedMessage, err := DecryptWithSymmetricKey(encryptedBytes, symmetricKey)
	if err != nil {
		debug.Log("server", fmt.Sprintf("Error decrypting message: %s", err.Error()))
		return "", err
	}
	return string(decryptedMessage), nil
}

func (network *Network) SendFirstMessage(peerIDs []string, receiver string) (models.FirstMessage, error) {
	sort.Strings(peerIDs)
	// Check if the user is online
	debug.Log("server", fmt.Sprintf("Checking if user %s is online", receiver))
	online := false
	for _, peer := range network.PubSubService.PeerList() {
		if peer.String() == receiver {
			online = true
			break
		}
	}
	if !online {
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

func (network *Network) runMonitoring(monitor *monitoring.SystemMonitor) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats, err := monitor.Collect()
			// if err == nil {
			// 	debug.Log("monitor", stats.String())
			// }
			statsMap := map[string]interface{}{
				"cpu": map[string]interface{}{
					"usage":        stats.CPUUsage,
					"processUsage": stats.ProcessCPU,
				},
				"memory": map[string]interface{}{
					"total":        stats.MemoryTotal,
					"used":         stats.MemoryUsed,
					"free":         stats.MemoryFree,
					"usagePercent": stats.MemoryUsagePerc,
					"processUsage": stats.ProcessMemory,
				},
				"network": map[string]interface{}{
					"bytesSent":       stats.BytesSent,
					"bytesReceived":   stats.BytesRecv,
					"packetsSent":     stats.PacketsSent,
					"packetsReceived": stats.PacketsRecv,
					"speed":           stats.NetworkSpeed,
				},
				"runtime": map[string]interface{}{
					"goroutines": stats.NumGoroutines,
					"numCPU":     stats.NumCPU,
				},
				"timestamp": stats.Timestamp,
			}
			jsonData, err := json.Marshal(statsMap)
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error marshaling stats map: %s", err.Error()))
			}
			// Append to the file
			os.MkdirAll("stats", 0755)
			file, err := os.OpenFile("stats/system_stats.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				debug.Log("server", fmt.Sprintf("Error opening file: %s", err.Error()))
			}
			defer file.Close()
			file.Write(jsonData)
			file.Write([]byte("\n"))
		}
	}
}
