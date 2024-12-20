package backend

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type Network struct {
	p2p      *P2P
	chatRoom *ChatRoom
}

func (network *Network) ConnectToNetwork() {
	// Define input flags
	username := "newuser"
	chatroom := "messagemesh"
	// Parse input flags

	// Set the log level
	logrus.SetLevel(logrus.InfoLevel)

	fmt.Println("The PeerChat Application is starting.")
	fmt.Println("This may take upto 30 seconds.")
	fmt.Println()

	// Create a new P2PHost
	network.p2p = NewP2P()
	logrus.Infoln("Completed P2P Setup")

	// Connect to peers with the chosen discovery method
	network.p2p.AdvertiseConnect()

	logrus.Infoln("Connected to Service Peers")

	// Join the chat room
	network.chatRoom, _ = JoinChatRoom(network.p2p, username, chatroom)
	logrus.Infof("Joined the '%s' chatroom as '%s'", network.chatRoom.RoomName, network.chatRoom.UserName)

	// Wait for network setup to complete
	time.Sleep(time.Second * 5)

	logrus.Infoln("Connected to Service Peers")

	// Print my peer ID
	logrus.Infof("My Peer ID: %s", network.chatRoom.SelfID())

	// Print my multiaddress
	logrus.Infof("My Multiaddress: %s", network.p2p.AllNodeAddr())

	// List of peers in the whole network
	// logrus.Infof("Connected to %s peers in the network", p2phost.Host.Peerstore().Peers())

	// List the dht peers
	// logrus.Infof("DHT Peers: %s", p2phost.KadDHT.RoutingTable().ListPeers())

	// Keep the main thread alive and check for new peers
	// networkPeerListCount := -1
	starteventhandler(network.chatRoom)

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
	// 	chatapp.Outbound <- fmt.Sprintf("Hello from %s", chatapp.UserName)
	// 	msg := <-chatapp.Inbound
	// 	logrus.Infof("Message: %s", msg.Message)
	// }
}

func starteventhandler(cr *ChatRoom) {
	refreshticker := time.NewTicker(time.Second)
	defer refreshticker.Stop()

	for {
		select {

		// case msg := <-cr.MsgInputs:
		// 	// Send the message to outbound queue
		// 	cr.Outbound <- msg
		// 	// Add the message to the message box as a self message
		// 	cr.display_selfmessage(msg)

		// case cmd := <-cr.CmdInputs:
		// 	// Handle the recieved command
		// 	go cr.handlecommand(cmd)

		case msg := <-cr.Inbound:
			// Print the recieved messages to the message box
			logrus.Infof("Message: %s", msg.Message)

		case log := <-cr.Logs:
			// Add the log to the message box
			logrus.Infof("Log: %s", log)

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
	network.chatRoom.Outbound <- message
}
