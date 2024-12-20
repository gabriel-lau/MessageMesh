package backend

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	blue   = "\033[34m"
	reset  = "\033[0m"
	purple = "\033[35m"
	pink   = "\033[95m"
)

type Network struct {
	P2p      *P2P
	ChatRoom *ChatRoom
}

func (network *Network) ConnectToNetwork() {
	// Define input flags
	username := "newuser"
	chatroom := "messagemesh"
	// Parse input flags

	fmt.Println(blue + "[server.go]" + reset + "The PeerChat Application is starting.")
	fmt.Println(blue + "[server.go]" + reset + "This may take upto 30 seconds.")

	// Create a new P2PHost
	network.P2p = NewP2P()
	fmt.Println(blue + "[server.go]" + reset + " Completed P2P Setup")

	// Connect to peers with the chosen discovery method
	network.P2p.AdvertiseConnect()

	fmt.Println(blue + "[server.go]" + reset + " Connected to Service Peers")

	// Join the chat room
	network.ChatRoom, _ = JoinChatRoom(network.P2p, username, chatroom)
	fmt.Printf(blue+"[server.go]"+reset+" Joined the '%s' chatroom as '%s'\n", network.ChatRoom.RoomName, network.ChatRoom.UserName)

	// Wait for network setup to complete
	time.Sleep(time.Second * 5)

	fmt.Println(blue + "[server.go]" + reset + " Connected to Service Peers")

	// Print my peer ID
	fmt.Printf(blue+"[server.go]"+reset+" My Peer ID: %s\n", network.ChatRoom.SelfID())

	// Print my multiaddress
	fmt.Printf(blue+"[server.go]"+reset+" My Multiaddress: %s\n", network.P2p.AllNodeAddr())

	// List of peers in the whole network
	// logrus.Infof("Connected to %s peers in the network", p2phost.Host.Peerstore().Peers())

	// List the dht peers
	// logrus.Infof("DHT Peers: %s", p2phost.KadDHT.RoutingTable().ListPeers())

	// Keep the main thread alive and check for new peers
	// networkPeerListCount := -1
	go network.starteventhandler()

	// channel <- &Network{
	// 	p2p:      p2phost,
	// 	chatRoom: chatapp,
	// }

	// chatRoomPeerListCount := -1
	// for {
	// 	// if len(p2phost.Host.Peerstore().Peers()) != networkPeerListCount {
	// 	// 	logrus.Infof("Connected to %d peers in the network", len(p2phost.Host.Peerstore().Peers()))
	// 	// 	networkPeerListCount = len(p2phost.Host.Peerstore().Peers())
	// 	// }

	// 	// Get the list of peers
	// 	if len(chatapp.PeerList()) != chatRoomPeerListCount {
	// 		logrus.Infof("Connected to %d peers in the chatroom", len(chatapp.PeerList()))
	// 		for _, p := range chatapp.PeerList() {
	// 			logrus.Infof("Peer ID: %s", p.String())
	// 		}
	// 		chatRoomPeerListCount = len(chatapp.PeerList())
	// 	}
	// 	chatapp.Outbound <- fmt.SPrintln("Hello from %s", chatapp.UserName)
	// 	msg := <-chatapp.Inbound
	// 	logrus.Infof("Message: %s", msg.Message)
	// }
}

func (network *Network) starteventhandler() {
	refreshticker := time.NewTicker(time.Second)
	defer refreshticker.Stop()

	chatRoomPeerListCount := -1
	for {
		if len(network.ChatRoom.PeerList()) != chatRoomPeerListCount {
			fmt.Printf(blue+"[server.go]"+reset+" Connected to %d peers in the chatroom\n", len(network.ChatRoom.PeerList()))
			for _, p := range network.ChatRoom.PeerList() {
				logrus.Infof("Peer ID: %s", p.String())
			}
			chatRoomPeerListCount = len(network.ChatRoom.PeerList())
		}
		select {

		// case msg := <-cr.MsgInputs:
		// 	// Send the message to outbound queue
		// 	cr.Outbound <- msg
		// 	// Add the message to the message box as a self message
		// 	cr.display_selfmessage(msg)

		// case cmd := <-cr.CmdInputs:
		// 	// Handle the recieved command
		// 	go cr.handlecommand(cmd)

		case msg := <-network.ChatRoom.Inbound:
			// Print the recieved messages to the message box
			fmt.Printf(blue+"[server.go]"+reset+" Message: %s\n", msg.Message)

		case log := <-network.ChatRoom.Logs:
			// Add the log to the message box
			fmt.Printf(blue+"[server.go]"+reset+" Log: %s\n", log)

			// case <-refreshticker.C:
			// 	// Refresh the list of peers in the chat room periodically
			// 	cr.syncpeerbox()

			// case <-cr.psctx.Done():
			// 	// End the event loop
			// 	return
		}
	}
}

func (network *Network) SendMessage(message string) {
	network.ChatRoom.Outbound <- message
}
