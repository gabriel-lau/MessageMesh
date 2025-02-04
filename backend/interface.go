package backend

import (
	"MessageMesh/backend/models"
	"context"

	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Network struct {
	P2p      *P2P
	ChatRoom *ChatRoom
}

type ChatRoom struct {

	// Represents the channel of incoming messages
	Inbound chan models.Message
	// Represents the channel of outgoing messages
	Outbound chan string
	// Represents the channel of chat log messages
	Logs chan chatlog

	PeerJoin  chan peer.ID
	PeerLeave chan peer.ID
	PeerIDs   chan []peer.ID

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

type P2P struct {
	// Represents the host context layer
	Ctx context.Context

	// Represents the libp2p host
	Host host.Host

	// Represents the DHT routing table
	KadDHT *dht.IpfsDHT

	// Represents the peer discovery service
	Discovery *discovery.RoutingDiscovery

	// Represents the PubSub Handler
	PubSub *pubsub.PubSub
}
