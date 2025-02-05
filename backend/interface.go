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
	// P2P Network (libp2p connections)
	P2p      *P2P
	ChatRoom *ChatRoom
}

type P2P struct {
	Ctx       context.Context
	Host      host.Host
	KadDHT    *dht.IpfsDHT
	Discovery *discovery.RoutingDiscovery
	PubSub    *pubsub.PubSub
}

type ChatRoom struct {
	Inbound   chan models.Message
	Outbound  chan string
	PeerJoin  chan peer.ID
	PeerLeave chan peer.ID
	PeerIDs   chan []peer.ID
	selfid    peer.ID
	psctx     context.Context
	pscancel  context.CancelFunc
	pstopic   *pubsub.Topic
	psub      *pubsub.Subscription
}
