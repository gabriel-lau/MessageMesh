package backend

import (
	"fmt"
	"time"
)

const (
	blue   = "\033[34m"
	reset  = "\033[0m"
	purple = "\033[35m"
	pink   = "\033[95m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
)

func (network *Network) ConnectToNetwork() {
	fmt.Println(blue + "[server.go]" + "[" + time.Now().Format("15:04:05") + "]" + reset + " The PeerChat Application is starting.")
	fmt.Println(blue + "[server.go]" + "[" + time.Now().Format("15:04:05") + "]" + reset + " This may take upto 30 seconds.")

	// Create a new P2PHost
	network.P2p = NewP2P()
	fmt.Println(blue + "[server.go]" + " [" + time.Now().Format("15:04:05") + "]" + reset + " Completed P2P Setup")

	// Connect to peers with the chosen discovery method
	network.P2p.AdvertiseConnect()
	// network.P2p.AnnounceConnect()
	fmt.Println(blue + "[server.go]" + " [" + time.Now().Format("15:04:05") + "]" + reset + " Connected to Service Peers")

	// Join the chat room
	network.ChatRoom, _ = JoinChatRoom(network.P2p, "username")
	fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Joined the '%s' chatroom as '%s'\n", network.ChatRoom.RoomName, network.ChatRoom.UserName)

	// Wait for network setup to complete
	time.Sleep(time.Second * 5)
	fmt.Println(blue + "[server.go]" + " [" + time.Now().Format("15:04:05") + "]" + reset + " Connected to Service Peers")

	// Print my peer ID
	fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" My Peer ID: %s\n", network.ChatRoom.SelfID())

	// Print my multiaddress
	fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" My Multiaddress: %s\n", network.P2p.AllNodeAddr())

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
			fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Message: %s\n", msg.Message)

		case log := <-network.ChatRoom.Logs:
			// Add the log to the message box
			fmt.Printf(blue+"[server.go]"+" ["+time.Now().Format("15:04:05")+"]"+reset+" Log: %s\n", log)

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
