package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func JoinPubSub(p2phost *P2PService) (*PubSubService, error) {

	// Create a PubSub topic with the room name
	topic, err := p2phost.PubSub.Join("messagemesh")
	// Check the error
	if err != nil {
		debug.Log("err", "Could not join the chat room")
		return nil, err
	}
	debug.Log("pubsub", "Joined the chat room")

	// Subscribe to the PubSub topic
	sub, err := topic.Subscribe()
	// Check the error
	if err != nil {
		debug.Log("err", "Could not subscribe to the chat room")
		return nil, err
	}
	debug.Log("pubsub", "Subscribed to the chat room")

	// Create cancellable context
	pubsubctx, cancel := context.WithCancel(context.Background())

	// Create a ChatRoom object
	pubsubservice := &PubSubService{
		Inbound:   make(chan models.Message),
		Outbound:  make(chan models.Message),
		PeerJoin:  make(chan peer.ID),
		PeerLeave: make(chan peer.ID),

		psctx:    pubsubctx,
		pscancel: cancel,
		pstopic:  topic,
		psub:     sub,
		selfid:   p2phost.Host.ID(),
	}

	// Start the subscribe loop
	go pubsubservice.SubLoop()
	debug.Log("pubsub", "SubLoop started")

	// Start the publish loop
	go pubsubservice.PubLoop()
	debug.Log("pubsub", "PubLoop started")

	// Start the peer joined loop
	go pubsubservice.PeerJoinedLoop()
	debug.Log("pubsub", "PeerJoinedLoop started")

	// Return the chatroom
	return pubsubservice, nil
}

// A method of ChatRoom that publishes a chatmessage
// to the PubSub topic until the pubsub context closes
func (pubSubService *PubSubService) PubLoop() {
	for {
		select {
		case <-pubSubService.psctx.Done():
			return

		case message := <-pubSubService.Outbound:
			// Create a ChatMessage

			// Marshal the ChatMessage into a JSON
			messagebytes, err := json.Marshal(message)
			if err != nil {
				debug.Log("err", "Could not marshal JSON")
				continue
			}
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Pub Message marshalled")

			// Publish the message to the topic
			err = pubSubService.pstopic.Publish(pubSubService.psctx, messagebytes)
			if err != nil {
				debug.Log("err", "Could not publish to topic")
				continue
			}
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Pub Message published")
		}
	}
}

// A method of ChatRoom that continously reads from the subscription
// until either the subscription or pubsub context closes.
// The recieved message is parsed sent into the inbound channel
func (pubSubService *PubSubService) SubLoop() {
	// Start loop
	for {
		select {
		case <-pubSubService.psctx.Done():
			return

		default:
			// Read a message from the subscription
			message, err := pubSubService.psub.Next(pubSubService.psctx)
			// Check error
			if err != nil {
				// Close the messages queue (subscription has closed)
				close(pubSubService.Inbound)
				debug.Log("err", "Subscription has closed")
				return
			}

			// Check if message is from self
			if message.ReceivedFrom == pubSubService.selfid {
				debug.Log("pubsub", "Sub Message from self")
			} else {
				debug.Log("pubsub", "Sub Message from other peer")
			}

			// Declare a ChatMessage
			cm := &models.Message{}
			// Unmarshal the message data into a ChatMessage
			err = json.Unmarshal(message.Data, cm)
			if err != nil {
				debug.Log("err", "Could not unmarshal JSON")
				continue
			}
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Sub Message unmarshalled")
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Sender: " + cm.Sender)
			// fmt.Println(green + "[chatRoom.go]" + " [" + time.Now().Format("15:04:05") + "] " + reset + "Receiver: " + cm.Receiver)

			// Send the ChatMessage into the message queue
			pubSubService.Inbound <- *cm
		}
	}
}

func (pubSubService *PubSubService) PeerJoinedLoop() {
	// Get the event handler for the topic
	evts, err := pubSubService.pstopic.EventHandler()
	if err != nil {
		debug.Log("err", fmt.Sprintf("Failed to get event handler: %s", err))
		return
	}

	for {
		peerEvent, err := evts.NextPeerEvent(context.Background())
		if err != nil {
			debug.Log("err", fmt.Sprintf("Failed to get next peer event: %s", err))
			continue
		}

		switch peerEvent.Type {
		case pubsub.PeerJoin: // PeerJoin event
			debug.Log("pubsub", fmt.Sprintf("Peer joined: %s", peerEvent.Peer))
			pubSubService.PeerJoin <- peerEvent.Peer
			// raftInstance.AddVoter(raft.ServerID(peerEvent.Peer.String()), raft.ServerAddress(peerEvent.Peer.String()), 0, 0)

		case pubsub.PeerLeave: // PeerLeave event
			debug.Log("pubsub", fmt.Sprintf("Peer left: %s", peerEvent.Peer))
			pubSubService.PeerLeave <- peerEvent.Peer
			// raftInstance.RemoveServer(raft.ServerID(peerEvent.Peer.String()), 0, 0)
		}
	}
}

// A method of ChatRoom that returns a list
// of all peer IDs connected to it
func (pubSubService *PubSubService) PeerList() []peer.ID {
	// Return the slice of peer IDs connected to chat room topic
	return pubSubService.pstopic.ListPeers()
}

// A method of ChatRoom that updates the chat
// room by subscribing to the new topic
func (pubSubService *PubSubService) Exit() {
	defer pubSubService.pscancel()

	// Cancel the existing subscription
	pubSubService.psub.Cancel()
	// Close the topic handler
	pubSubService.pstopic.Close()
}

// A method of ChatRoom that returns the self peer ID
func (pubSubService *PubSubService) SelfID() peer.ID {
	return pubSubService.selfid
}
