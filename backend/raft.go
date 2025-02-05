package backend

import (
	"MessageMesh/debug"
	"fmt"
	"io"
	"time"

	"MessageMesh/backend/models"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	consensus "github.com/libp2p/go-libp2p-consensus"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

type raftState struct {
	Blockchain Blockchain
}

type raftOP struct {
	Type    string // "ADD_MESSAGE_BLOCK" or "ADD_ACCOUNT_BLOCK"
	Message *models.Message
	Account *models.Account
}

func (o *raftOP) ApplyTo(state consensus.State) (consensus.State, error) {
	currentState := state.(*raftState)

	switch o.Type {
	case "ADD_MESSAGE_BLOCK":
		newBlock := currentState.Blockchain.AddMessageBlock(o.Message)
		debug.Log("blockchain", fmt.Sprintf("New message block added: %d", newBlock.Index))

	case "ADD_ACCOUNT_BLOCK":
		newBlock := currentState.Blockchain.AddAccountBlock(o.Account)
		debug.Log("blockchain", fmt.Sprintf("New account block added: %d", newBlock.Index))
	}

	return currentState, nil
}

func StartConsensus(network *Network) {
	// Initialize blockchain with genesis block
	initialState := &raftState{
		Blockchain: Blockchain{
			Chain: []*Block{CreateGenesisBlock()},
		},
	}

	// Create the consensus with blockchain state
	raftconsensus := libp2praft.NewOpLog(initialState, &raftOP{})

	pids := network.PubSubService.PeerList()
	// pids := network.P2p.Host.Peerstore().Peers()
	// -- Create Raft servers configuration
	pids = append(pids, network.P2pService.Host.ID())
	servers := make([]raft.Server, len(pids))
	for i, pid := range pids {
		servers[i] = raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(pid.String()),
			Address:  raft.ServerAddress(pid.String()),
		}
		debug.Log("raft", fmt.Sprintf("Server: %v", servers[i]))
	}
	serverConfig := raft.Configuration{
		Servers: servers,
	}
	// --

	// -- Create LibP2P transports Raft
	transport, err := libp2praft.NewLibp2pTransport(network.P2pService.Host, 3*time.Second)
	if err != nil {
		debug.Log("err", fmt.Sprintf("Failed to create LibP2P transport: %s", err))
	}
	// --

	// -- Configuration
	raftQuiet := false
	config := raft.DefaultConfig()
	if raftQuiet {
		config.LogOutput = io.Discard
		config.Logger = nil
	}
	config.LocalID = raft.ServerID(network.P2pService.Host.ID().String())
	config.HeartbeatTimeout = 1000 * time.Millisecond // Increase heartbeat timeout
	config.ElectionTimeout = 1000 * time.Millisecond  // Increase election timeout
	config.CommitTimeout = 500 * time.Millisecond     // Increase commit timeout
	config.LeaderLeaseTimeout = 1000 * time.Millisecond
	config.SnapshotInterval = 10 * time.Second
	// --

	// -- SnapshotStore
	var raftTmpFolder = "db/raft_testing_tmp"
	snapshots, err := raft.NewFileSnapshotStore(raftTmpFolder, 3, nil)
	if err != nil {
		debug.Log("err", fmt.Sprintf("Failed to create snapshot store: %s", err))
	}

	// -- Log store and stable store: we use inmem.
	// logStore := raft.NewInmemStore()
	logStore, _ := raftboltdb.NewBoltStore("db/raft.db")
	// --

	// -- Boostrap everything if necessary
	bootstrapped, err := raft.HasExistingState(logStore, logStore, snapshots)
	if err != nil {
		debug.Log("err", fmt.Sprintf("Failed to check existing state: %s", err))
	}

	// Only bootstrap if we're the first node (no peers)
	if !bootstrapped && len(pids) <= 1 {
		debug.Log("raft", "Bootstrapping new cluster")
		err := raft.BootstrapCluster(config, logStore, logStore, snapshots, transport, serverConfig)
		if err != nil {
			debug.Log("err", fmt.Sprintf("Failed to bootstrap cluster: %s", err))
		}
	} else {
		debug.Log("raft", "Joining existing cluster")
	}

	raftInstance, err := raft.NewRaft(config, raftconsensus.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		fmt.Println(err)
	}

	actor := libp2praft.NewActor(raftInstance)
	raftconsensus.SetActor(actor)
	consensusService := &ConsensusService{
		Raft:      raftInstance,
		Actor:     actor,
		Consensus: raftconsensus,
	}
	network.ConsensusService = consensusService

	go networkLoop(network, raftInstance)

	go blockchainLoop(network, raftInstance, raftconsensus, actor)
}

func networkLoop(network *Network, raftInstance *raft.Raft) {
	for {
		select {
		case peer := <-network.PubSubService.PeerJoin:
			debug.Log("raft", fmt.Sprintf("Peer joined: %s", peer))

			// Only add voter if we are the leader
			if raftInstance.State() == raft.Leader {
				future := raftInstance.AddVoter(
					raft.ServerID(peer.String()),
					raft.ServerAddress(peer.String()),
					0,
					5*time.Second,
				)
				if err := future.Error(); err != nil {
					debug.Log("err", fmt.Sprintf("Failed to add voter: %s", err))
				}
			}
			network.PubSubService.PeerIDs <- network.PubSubService.PeerList()

		case peer := <-network.PubSubService.PeerLeave:
			debug.Log("raft", fmt.Sprintf("Peer left: %s", peer))
			if raftInstance.State() == raft.Leader {
				future := raftInstance.RemoveServer(
					raft.ServerID(peer.String()),
					0,
					5*time.Second,
				)
				if err := future.Error(); err != nil {
					debug.Log("err", fmt.Sprintf("Failed to remove server: %s", err))
				}
			}
			network.PubSubService.PeerIDs <- network.PubSubService.PeerList()
		}
	}
}

func blockchainLoop(network *Network, raftInstance *raft.Raft, raftconsensus *libp2praft.Consensus, actor *libp2praft.Actor) {
	for {
		select {
		case <-raftconsensus.Subscribe():
			newState, _ := raftconsensus.GetCurrentState()
			blockchain := newState.(*raftState).Blockchain
			debug.Log("blockchain", fmt.Sprintf("Blockchain updated, current length: %d", len(blockchain.Chain)))
			latestBlock := blockchain.GetLatestBlock()

			// Type assertion to access specific data
			switch latestBlock.BlockType {
			case "message":
				if messageData, ok := latestBlock.Data.(*MessageData); ok {
					debug.Log("blockchain", fmt.Sprintf("Latest message from: %s", messageData.Message.Sender))
					debug.Log("blockchain", fmt.Sprintf("Latest message: %s", messageData.Message.Message))
				}
			case "account":
				if accountData, ok := latestBlock.Data.(*AccountData); ok {
					debug.Log("blockchain", fmt.Sprintf("Latest account: %s", accountData.Account.Username))
				}
			default:
				debug.Log("blockchain", fmt.Sprintf("Latest block type: %s", latestBlock.BlockType))
			}
			network.ConsensusService.Blockchain <- Block{
				Index:     latestBlock.Index,
				Timestamp: latestBlock.Timestamp,
				PrevHash:  latestBlock.PrevHash,
				Hash:      latestBlock.Hash,
				BlockType: latestBlock.BlockType,
				Data:      latestBlock.Data,
			}

		case <-raftInstance.LeaderCh():
			debug.Log("raft", "Leader changed")

		case <-network.PubSubService.Outbound:
			debug.Log("raft", "Outbound message received")

		case message := <-network.PubSubService.Inbound:
			debug.Log("blockchain", fmt.Sprintf("Received message: %s", message.Message))

			if actor.IsLeader() {
				// Create a message block using the new structure
				op := &raftOP{
					Type: "ADD_MESSAGE_BLOCK",
					Message: &models.Message{
						Sender:    message.Sender,
						Receiver:  message.Receiver,
						Message:   message.Message,
						Timestamp: time.Now().Format(time.RFC3339),
					},
				}

				_, err := raftconsensus.CommitOp(op)
				if err != nil {
					debug.Log("err", fmt.Sprintf("Failed to commit block: %s", err))
				}
			}
		}
	}
}

// func waitForLeader(r *raft.Raft) {
// 	obsCh := make(chan raft.Observation, 1)
// 	observer := raft.NewObserver(obsCh, false, nil)
// 	r.RegisterObserver(observer)
// 	defer r.DeregisterObserver(observer)

// 	// New Raft does not allow leader observation directy
// 	// What's worse, there will be no notification that a new
// 	// leader was elected because observations are set before
// 	// setting the Leader and only when the RaftState has changed.
// 	// Therefore, we need a ticker.

// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	ticker := time.NewTicker(time.Second / 2)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case obs := <-obsCh:
// 			switch obs.Data.(type) {
// 			case raft.RaftState:
// 				if r.Leader() != "" {
// 					return
// 				}
// 			}
// 		case <-ticker.C:
// 			if r.Leader() != "" {
// 				return
// 			}
// 		case <-ctx.Done():
// 			debug.Log("raft", "timed out waiting for Leader")
// 			debug.Log("raft", fmt.Sprintf("Current Raft State: %s", r.State()))
// 			debug.Log("raft", fmt.Sprintf("Current Leader: %s", r.Leader()))
// 			return
// 		}
// 	}
// }
