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
	// refreshticker := time.NewTicker(time.Second)
	// defer refreshticker.Stop()
	network.ChatRoom.Outbound <- "I am " + network.ChatRoom.selfid.String()
	for {
		select {

		// case msg := <-network.ChatRoom.MsgInputs:
		// 	// Send the message to outbound queue
		// 	network.ChatRoom.Outbound <- msg
		// 	// Add the message to the message box as a self message
		// 	fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Message: %s\n", msg)

		// case cmd := <-network.ChatRoom.CmdInputs:
		// 	// Handle the recieved command
		// 	fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Command: %s\n", cmd)

		case msg := <-network.ChatRoom.Inbound:
			// Print the recieved messages to the message box
			debug.Log("server", fmt.Sprintf("Message: %s", msg.Message))

		case log := <-network.ChatRoom.Logs:
			// Add the log to the message box
			debug.Log("server", fmt.Sprintf("Log: %s", log))

		// case <-refreshticker.C:
		// 	// Refresh the list of peers in the chat room periodically
		// 	peerlist := network.ChatRoom.PeerList()
		// 	fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Peers: %s\n", peerlist)
		// 	peerstoreList := network.P2p.Host.Peerstore().Peers()
		// 	fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Peerstore: %s\n", peerstoreList)

		case <-network.ChatRoom.psctx.Done():
			// End the event loop
			return
		}
	}
}

func (network *Network) SendMessage(message string) {
	network.ChatRoom.Outbound <- message
}
