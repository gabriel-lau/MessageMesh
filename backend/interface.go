package backend

import (
	"MessageMesh/backend/models"
	"context"

	"github.com/hashicorp/raft"
	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

type Network struct {
	// P2P Service (libp2p connections)
	P2pService *P2PService
	// PubSub Service (PubSub connections)
	PubSubService *PubSubService
	// Consensus Service (Raft consensus)
	ConsensusService *ConsensusService
}

type P2PService struct {
	// Context
	Ctx context.Context
	// Host
	Host host.Host
	// KadDHT
	KadDHT *dht.IpfsDHT
	// Discovery
	Discovery *discovery.RoutingDiscovery
	// PubSub
	PubSub *pubsub.PubSub
}

type PubSubService struct {
	// Topic
	Topic string
	// Listen to new messages
	Inbound chan any
	// Send messages
	Outbound chan any
	// Listen to new peers
	PeerJoin chan peer.ID
	// Listen to peer leave
	PeerLeave chan peer.ID
	// List of peer IDs
	PeerIDs chan []peer.ID
	// Self peer ID
	selfid peer.ID
	// Context
	psctx context.Context
	// Cancel context
	pscancel context.CancelFunc
	// PubSub topic
	pstopic *pubsub.Topic
	// PubSub subscription
	psub *pubsub.Subscription
}

type ConsensusService struct {
	// Listen to latest block in blockchain
	LatestBlock chan models.Block
	// Blockchain
	Blockchain *models.Blockchain
	// Consensus connected
	Connected chan bool
	// Raft instance
	Raft *raft.Raft
	// Libp2p Raft actor
	Actor *libp2praft.Actor
	// Libp2p Raft consensus
	Consensus *libp2praft.Consensus
}
