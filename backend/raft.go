package backend

import (
	"MessageMesh/debug"
	"context"
	"fmt"
	"io"
	"time"

	"MessageMesh/backend/models"

	"github.com/hashicorp/raft"
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

func StartRaft(network *Network) {
	// Initialize blockchain with genesis block
	initialState := &raftState{
		Blockchain: Blockchain{
			Chain: []Block{CreateGenesisBlock()},
		},
	}

	// Create the consensus with blockchain state
	raftconsensus := libp2praft.NewOpLog(initialState, &raftOP{})

	pids := network.ChatRoom.PeerList()
	// pids := network.P2p.Host.Peerstore().Peers()
	// -- Create Raft servers configuration
	pids = append(pids, network.P2p.Host.ID())
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
	transport, err := libp2praft.NewLibp2pTransport(network.P2p.Host, 3*time.Second)
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
	config.LocalID = raft.ServerID(network.P2p.Host.ID().String())
	// --

	// -- SnapshotStore
	var raftTmpFolder = "db/raft_testing_tmp"
	snapshots, err := raft.NewFileSnapshotStore(raftTmpFolder, 3, nil)
	if err != nil {
		debug.Log("err", fmt.Sprintf("Failed to create snapshot store: %s", err))
	}

	// -- Log store and stable store: we use inmem.
	logStore := raft.NewInmemStore()
	// logStore, _ := raftboltdb.NewBoltStore("db/raft.db")
	// --

	// -- Boostrap everything if necessary
	bootstrapped, err := raft.HasExistingState(logStore, logStore, snapshots)
	if err != nil {
		debug.Log("err", fmt.Sprintf("Failed to check existing state: %s", err))
	}

	if !bootstrapped {
		// Bootstrap cluster first
		raft.BootstrapCluster(config, logStore, logStore, snapshots, transport, serverConfig)
	} else {
		debug.Log("raft", "Already initialized!")
	}

	raftInstance, err := raft.NewRaft(config, raftconsensus.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		fmt.Println(err)
	}

	actor := libp2praft.NewActor(raftInstance)
	raftconsensus.SetActor(actor)

	// waitForLeader(raftInstance)

	go networkLoop(network, raftInstance)

	go blockchainLoop(network, raftInstance, raftconsensus, actor)
}

// func updateState(c *libp2praft.Consensus) {
// 	loc, _ := time.LoadLocation("UTC")
// 	newState := &raftState{Now: time.Now().In(loc).String()}
// 	agreedState, err := c.CommitState(newState)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	if agreedState == nil {
// 		fmt.Println("agreedState is nil: commited on a non-leader?")
// 	}
// }

// func getState(c *libp2praft.Consensus) {
// 	state, err := c.GetCurrentState()
// 	if err != nil {
// 		debug.Log("err", err.Error())
// 	}
// 	if state == nil {
// 		debug.Log("err", "state is nil: commited on a non-leader?")
// 		return
// 	}
// 	debug.Log("raft", fmt.Sprintf("Current state: %s", state.(*raftState).Blockchain.GetLatestBlock().BlockType))
// }

func networkLoop(network *Network, raftInstance *raft.Raft) {
	// Listen for peer joins and leaves
	for {
		select {
		case peer := <-network.ChatRoom.PeerJoin:
			debug.Log("raft", fmt.Sprintf("Peer joined: %s", peer))
			raftInstance.AddVoter(raft.ServerID(peer.String()), raft.ServerAddress(peer.String()), 0, 0)
		case peer := <-network.ChatRoom.PeerLeave:
			debug.Log("raft", fmt.Sprintf("Peer left: %s", peer))
			raftInstance.RemoveServer(raft.ServerID(peer.String()), 0, 0)
			if raftInstance.Leader() == raft.ServerAddress(peer.String()) {
			}
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

		case <-raftInstance.LeaderCh():
			debug.Log("raft", "Leader changed")

		case <-network.ChatRoom.Outbound:
			debug.Log("raft", "Outbound message received")

		case message := <-network.ChatRoom.Inbound:
			debug.Log("blockchain", fmt.Sprintf("Received message: %s", message.Message))

			if actor.IsLeader() {
				// Example of creating a message block
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

func waitForLeader(r *raft.Raft) {
	obsCh := make(chan raft.Observation, 1)
	observer := raft.NewObserver(obsCh, false, nil)
	r.RegisterObserver(observer)
	defer r.DeregisterObserver(observer)

	// New Raft does not allow leader observation directy
	// What's worse, there will be no notification that a new
	// leader was elected because observations are set before
	// setting the Leader and only when the RaftState has changed.
	// Therefore, we need a ticker.

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	for {
		select {
		case obs := <-obsCh:
			switch obs.Data.(type) {
			case raft.RaftState:
				if r.Leader() != "" {
					return
				}
			}
		case <-ticker.C:
			if r.Leader() != "" {
				return
			}
		case <-ctx.Done():
			debug.Log("raft", "timed out waiting for Leader")
			debug.Log("raft", fmt.Sprintf("Current Raft State: %s", r.State()))
			debug.Log("raft", fmt.Sprintf("Current Leader: %s", r.Leader()))
			return
		}
	}
}
