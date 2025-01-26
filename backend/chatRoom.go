package backend

import (
	"MessageMesh/backend/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/raft"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// A structure that represents a PubSub Chat Room
// A structure that represents a chat log
type chatlog struct {
	logprefix string
	logmsg    string
}

// A constructor function that generates and returns a new
// ChatRoom for a given P2PHost, username and roomname
func JoinChatRoom(p2phost *P2P, username string) (*ChatRoom, error) {

	// Create a PubSub topic with the room name
	topic, err := p2phost.PubSub.Join("messagemesh")
	// Check the error
	if err != nil {
		fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Could not join the chat room")
		return nil, err
	}
	fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Joined the chat room")

	// Subscribe to the PubSub topic
	sub, err := topic.Subscribe()
	// Check the error
	if err != nil {
		fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Could not subscribe to the chat room")
		return nil, err
	}
	fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Subscribed to the chat room")

	// Create cancellable context
	pubsubctx, cancel := context.WithCancel(context.Background())

	// Create a ChatRoom object
	chatroom := &ChatRoom{
		Inbound:  make(chan models.Message),
		Outbound: make(chan string),
		Logs:     make(chan chatlog),

		psctx:    pubsubctx,
		pscancel: cancel,
		pstopic:  topic,
		psub:     sub,

		RoomName: "messagemesh",
		UserName: username,
		selfid:   p2phost.Host.ID(),
	}

	// Start the subscribe loop
	go chatroom.SubLoop()
	fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "SubLoop started")

	// Start the publish loop
	go chatroom.PubLoop()
	fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "PubLoop started")

	// Return the chatroom
	return chatroom, nil
}

// A method of ChatRoom that publishes a chatmessage
// to the PubSub topic until the pubsub context closes
func (cr *ChatRoom) PubLoop() {
	for {
		select {
		case <-cr.psctx.Done():
			return

		case message := <-cr.Outbound:
			// Create a ChatMessage
			m := models.Message{
				Sender:    cr.selfid.Pretty(),
				Receiver:  "QmYvjPHjCwsMXQThevzPyHTWwBK7VLHaAwjocEa42CK2vQ",
				Message:   message,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			}

			// Marshal the ChatMessage into a JSON
			messagebytes, err := json.Marshal(m)
			if err != nil {
				cr.Logs <- chatlog{logprefix: "puberr", logmsg: "could not marshal JSON"}
				fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Could not marshal JSON")
				continue
			}
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Pub Message marshalled")

			// Publish the message to the topic
			err = cr.pstopic.Publish(cr.psctx, messagebytes)
			if err != nil {
				cr.Logs <- chatlog{logprefix: "puberr", logmsg: "could not publish to topic"}
				fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Could not publish to topic")
				continue
			}
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Pub Message published")
		}
	}
}

// A method of ChatRoom that continously reads from the subscription
// until either the subscription or pubsub context closes.
// The recieved message is parsed sent into the inbound channel
func (cr *ChatRoom) SubLoop() {
	// Start loop
	for {
		select {
		case <-cr.psctx.Done():
			return

		default:
			// Read a message from the subscription
			message, err := cr.psub.Next(cr.psctx)
			// Check error
			if err != nil {
				// Close the messages queue (subscription has closed)
				close(cr.Inbound)
				cr.Logs <- chatlog{logprefix: "suberr", logmsg: "subscription has closed"}
				fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Subscription has closed")
				return
			}

			// Check if message is from self
			if message.ReceivedFrom == cr.selfid {
				fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Sub Message from self")
				continue
			} else {
				fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Sub Message from other peer")
			}

			// Declare a ChatMessage
			cm := &models.Message{}
			// Unmarshal the message data into a ChatMessage
			err = json.Unmarshal(message.Data, cm)
			if err != nil {
				cr.Logs <- chatlog{logprefix: "suberr", logmsg: "could not unmarshal JSON"}
				// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Could not unmarshal JSON")
				continue
			}
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Sub Message unmarshalled")
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Sender: " + cm.Sender)
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Receiver: " + cm.Receiver)

			// Send the ChatMessage into the message queue
			cr.Inbound <- *cm
		}
	}
}

func (cr *ChatRoom) PeerJoinedLoop(raftInstance *raft.Raft) {
	// Get the event handler for the topic
	evts, err := cr.pstopic.EventHandler()
	if err != nil {
		fmt.Println(green+"[chatRoom.go]"+" ["+time.Now().Format("15:04:05")+"] "+reset+"Failed to get event handler:", err)
		return
	}

	for {
		peerEvent, err := evts.NextPeerEvent(context.Background())
		if err != nil {
			fmt.Println(green+"[chatRoom.go]"+" ["+time.Now().Format("15:04:05")+"] "+reset+"Failed to get next peer event:", err)
			continue
		}

		switch peerEvent.Type {
		case pubsub.PeerJoin: // PeerJoin event
			fmt.Printf(green+"[chatRoom.go]"+" ["+time.Now().Format("15:04:05")+"] "+reset+"Peer joined: %s\n", peerEvent.Peer)
			raftInstance.AddVoter(raft.ServerID(peerEvent.Peer.String()), raft.ServerAddress(peerEvent.Peer.String()), 0, 0)

		case pubsub.PeerLeave: // PeerLeave event
			fmt.Printf(green+"[chatRoom.go]"+" ["+time.Now().Format("15:04:05")+"] "+reset+"Peer left: %s\n", peerEvent.Peer)
			raftInstance.RemoveServer(raft.ServerID(peerEvent.Peer.String()), 0, 0)
		}
	}
}

// A method of ChatRoom that returns a list
// of all peer IDs connected to it
func (cr *ChatRoom) PeerList() []peer.ID {
	// Return the slice of peer IDs connected to chat room topic
	return cr.pstopic.ListPeers()
}

// A method of ChatRoom that updates the chat
// room by subscribing to the new topic
func (cr *ChatRoom) Exit() {
	defer cr.pscancel()

	// Cancel the existing subscription
	cr.psub.Cancel()
	// Close the topic handler
	cr.pstopic.Close()
}

// A method of ChatRoom that returns the self peer ID
func (cr *ChatRoom) SelfID() peer.ID {
	return cr.selfid
}
