package backend

import (
	"MessageMesh/debug"
	"fmt"
	"time"
)

func (network *Network) ConnectToNetwork() {
	debug.Log("server", "The PeerChat Application is starting.")
	debug.Log("server", "This may take upto 30 seconds.")

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

	// Print my peer ID
	debug.Log("server", fmt.Sprintf("My Peer ID: %s", network.PubSubService.SelfID()))

	// Print my multiaddress
	debug.Log("server", fmt.Sprintf("My Multiaddress: %s", network.P2pService.AllNodeAddr()))

	go StartRaft(network)
}

func (network *Network) SendMessage(message string) {
	network.PubSubService.Outbound <- message
}
