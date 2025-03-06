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

// Define a message envelope structure
type MessageEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

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
		Inbound:   make(chan any),
		Outbound:  make(chan any),
		PeerJoin:  make(chan peer.ID, 10),
		PeerLeave: make(chan peer.ID, 10),
		PeerIDs:   make(chan []peer.ID, 10),
		psctx:     pubsubctx,
		pscancel:  cancel,
		pstopic:   topic,
		psub:      sub,
		selfid:    p2phost.Host.ID(),
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

		case packet := <-pubSubService.Outbound:
			// Create a ChatMessage

			// Marshal the ChatMessage into a JSON
			messagebytes, err := json.Marshal(packet)
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
			packet, err := pubSubService.psub.Next(pubSubService.psctx)
			// Check error
			if err != nil {
				// Close the messages queue (subscription has closed)
				close(pubSubService.Inbound)
				debug.Log("err", "Subscription has closed")
				return
			}

			// Check if message is from self
			if packet.ReceivedFrom == pubSubService.selfid {
				debug.Log("pubsub", "Sub Message from self")
			} else {
				debug.Log("pubsub", "Sub Message from other peer")
			}

			// First unmarshal just the envelope to determine the message type
			envelope := &MessageEnvelope{}
			err = json.Unmarshal(packet.Data, envelope)
			if err != nil {
				debug.Log("err", "Could not unmarshal message envelope: "+err.Error())
				continue
			}

			// Based on the type, unmarshal into the appropriate struct
			switch envelope.Type {
			case "Message":
				message := &models.Message{}
				if err := json.Unmarshal(envelope.Data, message); err != nil {
					debug.Log("err", "Could not unmarshal Message: "+err.Error())
					continue
				}
				pubSubService.Inbound <- *message

			case "FirstMessage":
				firstMessage := &models.FirstMessage{}
				if err := json.Unmarshal(envelope.Data, firstMessage); err != nil {
					debug.Log("err", "Could not unmarshal FirstMessage: "+err.Error())
					continue
				}
				pubSubService.Inbound <- *firstMessage

			case "Account":
				account := &models.Account{}
				if err := json.Unmarshal(envelope.Data, account); err != nil {
					debug.Log("err", "Could not unmarshal Account: "+err.Error())
					continue
				}
				pubSubService.Inbound <- *account

			default:
				debug.Log("warn", "Unknown message type: "+envelope.Type)
			}
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

	// Initialize PeerIDs channel if not already initialized
	if pubSubService.PeerIDs == nil {
		pubSubService.PeerIDs = make(chan []peer.ID, 10)
	}

	for {
		peerEvent, err := evts.NextPeerEvent(context.Background())
		if err != nil {
			debug.Log("err", fmt.Sprintf("Failed to get next peer event: %s", err))
			continue
		}

		switch peerEvent.Type {
		case pubsub.PeerJoin:
			debug.Log("pubsub", fmt.Sprintf("Peer joined: %s", peerEvent.Peer))
			select {
			case pubSubService.PeerJoin <- peerEvent.Peer:
				debug.Log("pubsub", fmt.Sprintf("Successfully sent peer join event for: %s", peerEvent.Peer))
			default:
				debug.Log("err", fmt.Sprintf("Channel blocked, couldn't send peer join event for: %s", peerEvent.Peer))
			}

			// Send updated peer list
			select {
			case pubSubService.PeerIDs <- pubSubService.PeerList():
				debug.Log("pubsub", "Sent updated peer list")
			default:
				debug.Log("err", "Channel blocked, couldn't send updated peer list")
			}

		case pubsub.PeerLeave:
			debug.Log("pubsub", fmt.Sprintf("Peer left: %s", peerEvent.Peer))
			select {
			case pubSubService.PeerLeave <- peerEvent.Peer:
				debug.Log("pubsub", fmt.Sprintf("Successfully sent peer leave event for: %s", peerEvent.Peer))
			default:
				debug.Log("err", fmt.Sprintf("Channel blocked, couldn't send peer leave event for: %s", peerEvent.Peer))
			}

			// Send updated peer list
			select {
			case pubSubService.PeerIDs <- pubSubService.PeerList():
				debug.Log("pubsub", "Sent updated peer list")
			default:
				debug.Log("err", "Channel blocked, couldn't send updated peer list")
			}
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
