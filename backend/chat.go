package backend

import (
	"context"
	"encoding/json"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Represents the default fallback room and user names
// if they aren't provided when the app is started
const defaultuser = "newuser"
const defaultroom = "lobby"

// A structure that represents a PubSub Chat Room
type ChatRoom struct {
	// Represents the P2P Host for the ChatRoom
	Host *P2P

	// Represents the channel of incoming messages
	Inbound chan chatmessage
	// Represents the channel of outgoing messages
	Outbound chan string
	// Represents the channel of chat log messages
	Logs chan chatlog

	// Represents the name of the chat room
	RoomName string
	// Represent the name of the user in the chat room
	UserName string
	// Represents the host ID of the peer
	selfid peer.ID

	// Represents the chat room lifecycle context
	psctx context.Context
	// Represents the chat room lifecycle cancellation function
	pscancel context.CancelFunc
	// Represents the PubSub Topic of the ChatRoom
	pstopic *pubsub.Topic
	// Represents the PubSub Subscription for the topic
	psub *pubsub.Subscription
}

// A structure that represents a chat message
type chatmessage struct {
	Message    string `json:"message"`
	SenderID   string `json:"senderid"`
	SenderName string `json:"sendername"`
}

// A structure that represents a chat log
type chatlog struct {
	logprefix string
	logmsg    string
}

// A constructor function that generates and returns a new
// ChatRoom for a given P2PHost, username and roomname
func JoinChatRoom(p2phost *P2P, username string, roomname string) (*ChatRoom, error) {

	// Create a PubSub topic with the room name
	topic, err := p2phost.PubSub.Join(roomname)
	// Check the error
	if err != nil {
		fmt.Println(pink + "[chat.go]" + reset + "Could not join the chat room")
		return nil, err
	}
	fmt.Println(pink + "[chat.go]" + reset + "Joined the chat room")

	// Subscribe to the PubSub topic
	sub, err := topic.Subscribe()
	// Check the error
	if err != nil {
		fmt.Println(pink + "[chat.go]" + reset + "Could not subscribe to the chat room")
		return nil, err
	}
	fmt.Println(pink + "[chat.go]" + reset + "Subscribed to the chat room")

	// Check the provided username
	if username == "" {
		// Use the default user name
		username = defaultuser
	}

	// Check the provided roomname
	if roomname == "" {
		// Use the default room name
		roomname = defaultroom
	}

	// Create cancellable context
	pubsubctx, cancel := context.WithCancel(context.Background())

	// Create a ChatRoom object
	chatroom := &ChatRoom{
		Host: p2phost,

		Inbound:  make(chan chatmessage),
		Outbound: make(chan string),
		Logs:     make(chan chatlog),

		psctx:    pubsubctx,
		pscancel: cancel,
		pstopic:  topic,
		psub:     sub,

		RoomName: roomname,
		UserName: username,
		selfid:   p2phost.Host.ID(),
	}

	// Start the subscribe loop
	go chatroom.SubLoop()
	fmt.Println(pink + "[chat.go]" + reset + "SubLoop started")

	// Start the publish loop
	go chatroom.PubLoop()
	fmt.Println(pink + "[chat.go]" + reset + "PubLoop started")

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
			m := chatmessage{
				Message:    message,
				SenderID:   cr.selfid.String(),
				SenderName: cr.UserName,
			}

			// Marshal the ChatMessage into a JSON
			messagebytes, err := json.Marshal(m)
			if err != nil {
				cr.Logs <- chatlog{logprefix: "puberr", logmsg: "could not marshal JSON"}
				fmt.Println(pink + "[chat.go]" + reset + "Could not marshal JSON")
				continue
			}
			fmt.Println(pink + "[chat.go]" + reset + "Message marshalled")

			// Publish the message to the topic
			err = cr.pstopic.Publish(cr.psctx, messagebytes)
			if err != nil {
				cr.Logs <- chatlog{logprefix: "puberr", logmsg: "could not publish to topic"}
				fmt.Println(pink + "[chat.go]" + reset + "Could not publish to topic")
				continue
			}
			fmt.Println(pink + "[chat.go]" + reset + "Message published")
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
				fmt.Println(pink + "[chat.go]" + reset + "Subscription has closed")
				return
			}
			fmt.Println(pink + "[chat.go]" + reset + "Message recieved")

			// Check if message is from self
			if message.ReceivedFrom == cr.selfid {
				continue
			}

			// Declare a ChatMessage
			cm := &chatmessage{}
			// Unmarshal the message data into a ChatMessage
			err = json.Unmarshal(message.Data, cm)
			if err != nil {
				cr.Logs <- chatlog{logprefix: "suberr", logmsg: "could not unmarshal JSON"}
				fmt.Println(pink + "[chat.go]" + reset + "Could not unmarshal JSON")
				continue
			}
			fmt.Println(pink + "[chat.go]" + reset + "Message unmarshalled")

			// Send the ChatMessage into the message queue
			cr.Inbound <- *cm
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

// A method of ChatRoom that updates the chat user name
func (cr *ChatRoom) UpdateUser(username string) {
	cr.UserName = username
}

// A method of ChatRoom that returns the self peer ID
func (cr *ChatRoom) SelfID() peer.ID {
	return cr.selfid
}
