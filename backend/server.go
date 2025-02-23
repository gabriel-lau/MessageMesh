package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"fmt"
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
	network.PubSubService.Outbound <- models.Message{
		Sender:    network.PubSubService.SelfID().String(),
		Receiver:  receiver,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
