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
	// Progress
	Progress *Progress
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
	// Raft instance
	Raft *raft.Raft
	// Libp2p Raft actor
	Actor *libp2praft.Actor
	// Libp2p Raft consensus
	Consensus *libp2praft.Consensus
}

type Progress struct {
	NetworkConnected chan bool
	PubSubJoined     chan bool
	BlockchainLoaded chan bool
}

// NewProgress creates a new Progress tracker with buffered channels
func NewProgress() *Progress {
	return &Progress{
		NetworkConnected: make(chan bool, 1),
		PubSubJoined:     make(chan bool, 1),
		BlockchainLoaded: make(chan bool, 1),
	}
}

// GetCurrentStep returns the current step in the startup process
// This should be called in a separate goroutine
func (p *Progress) GetCurrentStep() string {
	select {
	case <-p.NetworkConnected:
		select {
		case <-p.PubSubJoined:
			select {
			case <-p.BlockchainLoaded:
				return "Application ready"
			default:
				return "Loading blockchain data..."
			}
		default:
			return "Joining peer-to-peer network..."
		}
	default:
		return "Connecting to network..."
	}
}

// CompleteStep marks a step as complete by sending true to its channel
func (p *Progress) CompleteStep(step string) {
	switch step {
	case "network":
		p.NetworkConnected <- true
	case "pubsub":
		p.PubSubJoined <- true
	case "blockchain":
		p.BlockchainLoaded <- true
	}
}

// WaitForCompletion blocks until all steps are completed
func (p *Progress) WaitForCompletion() {
	<-p.NetworkConnected
	<-p.PubSubJoined
	<-p.BlockchainLoaded
}
