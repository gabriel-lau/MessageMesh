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
	network.P2p = NewP2P()
	debug.Log("server", "Completed P2P Setup")

	// Connect to peers with the chosen discovery method
	network.P2p.AdvertiseConnect()
	// network.P2p.AnnounceConnect()
	debug.Log("server", "Connected to Service Peers")

	// Join the chat room
	network.ChatRoom, _ = JoinChatRoom(network.P2p, "username")
	debug.Log("server", fmt.Sprintf("Joined the '%s' chatroom as '%s'", network.ChatRoom.RoomName, network.ChatRoom.UserName))

	// Wait for network setup to complete
	time.Sleep(time.Second * 5)
	debug.Log("server", "Connected to Service Peers")

	// Print my peer ID
	debug.Log("server", fmt.Sprintf("My Peer ID: %s", network.ChatRoom.SelfID()))

	// Print my multiaddress
	debug.Log("server", fmt.Sprintf("My Multiaddress: %s", network.P2p.AllNodeAddr()))

	go network.starteventhandler()

	go StartRaft(network)
}

func (network *Network) starteventhandler() {
	network.ChatRoom.Outbound <- "I am " + network.ChatRoom.selfid.String()
	// for {
	// 	select {

	// 	case msg := <-network.ChatRoom.Inbound:
	// 		// Print the recieved messages to the message box
	// 		debug.Log("server", fmt.Sprintf("Message: %s", msg.Message))

	// 	case log := <-network.ChatRoom.Logs:
	// 		// Add the log to the message box
	// 		debug.Log("server", fmt.Sprintf("Log: %s", log))

	// 	case <-network.ChatRoom.psctx.Done():
	// 		// End the event loop
	// 		return
	// 	}
	// }
}

func (network *Network) SendMessage(message string) {
	network.ChatRoom.Outbound <- message
}
