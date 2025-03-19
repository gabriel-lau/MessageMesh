package backend

import (
	"MessageMesh/debug"
	"fmt"
	"sort"
	"time"

	"MessageMesh/backend/models"

	"github.com/hashicorp/raft"
	consensus "github.com/libp2p/go-libp2p-consensus"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

type raftState struct {
	Blockchain models.Blockchain
}

type raftOP struct {
	Type         string // "ADD_MESSAGE_BLOCK" or "ADD_ACCOUNT_BLOCK" or "ADD_FIRST_MESSAGE_BLOCK"
	Message      *models.Message
	Account      *models.Account
	FirstMessage *models.FirstMessage
}

func (o *raftOP) ApplyTo(state consensus.State) (consensus.State, error) {
	currentState := state.(*raftState)

	switch o.Type {
	case "ADD_MESSAGE_BLOCK":
		if o.Message.Sender == "" || o.Message.Receiver == "" || o.Message.Message == "" {
			return currentState, fmt.Errorf("message is missing required fields")
		}
		if o.Message.Sender == o.Message.Receiver {
			return currentState, fmt.Errorf("message sender and receiver cannot be the same")
		}
	case "ADD_ACCOUNT_BLOCK":
		if o.Account.Username == "" {
			return currentState, fmt.Errorf("account is missing required fields")
		}
	case "ADD_FIRST_MESSAGE_BLOCK":
		if len(o.FirstMessage.PeerIDs) != 2 {
			return currentState, fmt.Errorf("first message must have exactly 2 peer IDs")
		}
		if o.FirstMessage.PeerIDs[0] == o.FirstMessage.PeerIDs[1] {
			return currentState, fmt.Errorf("first message peer IDs cannot be the same")
		}
		if o.FirstMessage.PeerIDs[0] == "" || o.FirstMessage.PeerIDs[1] == "" {
			return currentState, fmt.Errorf("first message peer IDs cannot be empty")
		}
		if o.FirstMessage.SymetricKey0 == nil || o.FirstMessage.SymetricKey1 == nil {
			return currentState, fmt.Errorf("first message symetric keys cannot be empty")
		}
	}

	// Apply the operation if validation passed
	switch o.Type {
	case "ADD_MESSAGE_BLOCK":
		newBlock := currentState.Blockchain.AddMessageBlock(*o.Message)
		debug.Log("raft", fmt.Sprintf("New message block added: %d", newBlock.Index))

	case "ADD_ACCOUNT_BLOCK":
		newBlock := currentState.Blockchain.AddAccountBlock(*o.Account)
		debug.Log("raft", fmt.Sprintf("New account block added: %d", newBlock.Index))

	case "ADD_FIRST_MESSAGE_BLOCK":
		newBlock := currentState.Blockchain.AddFirstMessageBlock(*o.FirstMessage)
		debug.Log("raft", fmt.Sprintf("New first message block added: %d", newBlock.Index))
	}

	return currentState, nil
}

func StartConsensus(network *Network) (*ConsensusService, error) {
	// Initialize blockchain with genesis block
	initialState := &raftState{
		Blockchain: models.Blockchain{
			Chain: []*models.Block{models.CreateGenesisBlock()},
		},
	}

	// Create the consensus with blockchain state
	raftconsensus := libp2praft.NewOpLog(initialState, &raftOP{})

	pids := network.PubSubService.PeerList()
	pids = append(pids, network.P2pService.Host.ID())
	servers := make([]raft.Server, len(pids))
	for i, pid := range pids {
		servers[i] = raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(pid.String()),
			Address:  raft.ServerAddress(pid.String()),
		}
	}
	serverConfig := raft.Configuration{
		Servers: servers,
	}

	transport, err := libp2praft.NewLibp2pTransport(network.P2pService.Host, 3*time.Second)
	if err != nil {
		return nil, err
	}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(network.P2pService.Host.ID().String())
	config.HeartbeatTimeout = 1000 * time.Millisecond
	config.ElectionTimeout = 1000 * time.Millisecond
	config.CommitTimeout = 500 * time.Millisecond
	config.LeaderLeaseTimeout = 1000 * time.Millisecond

	snapshots, err := raft.NewFileSnapshotStore("db/raft_testing_tmp", 3, nil)
	if err != nil {
		return nil, err
	}

	logStore := raft.NewInmemStore()
	// logStore, _ := raftboltdb.NewBoltStore("db/raft.db")

	// Check if we're the first node
	isFirstNode := len(pids) <= 1

	// Only bootstrap if we're the first node
	if isFirstNode {
		debug.Log("raft", "Bootstrapping new cluster as first node")
		if err := raft.BootstrapCluster(config, logStore, logStore, snapshots, transport, serverConfig); err != nil {
			return nil, fmt.Errorf("bootstrap error: %v", err)
		}
	}

	raftInstance, err := raft.NewRaft(config, raftconsensus.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	// If we're not the first node, wait for the leader to add us
	if !isFirstNode {
		debug.Log("raft", "Waiting to be added to existing cluster...")
		// The leader will add us through the networkLoop
	}

	actor := libp2praft.NewActor(raftInstance)
	raftconsensus.SetActor(actor)

	consensusService := &ConsensusService{
		LatestBlock: make(chan models.Block),
		Blockchain:  &initialState.Blockchain,
		Connected:   make(chan bool),
		Raft:        raftInstance,
		Actor:       actor,
		Consensus:   raftconsensus,
	}

	go networkLoop(network, raftInstance)
	go blockchainLoop(network, raftInstance, raftconsensus, actor)

	return consensusService, nil
}

func networkLoop(network *Network, raftInstance *raft.Raft) {
	debug.Log("raft", "Starting network loop")

	// Verify channel connection
	debug.Log("raft", "Checking PubSubService channels...")
	if network.PubSubService == nil {
		debug.Log("raft", "ERROR: PubSubService is nil")
		return
	}
	if network.PubSubService.PeerJoin == nil {
		debug.Log("raft", "ERROR: PeerJoin channel is nil")
		return
	}

	debug.Log("raft", "Channels verified, starting main loop")

	for {
		select {
		case peer := <-network.PubSubService.PeerJoin:
			debug.Log("raft", fmt.Sprintf("Network loop received peer join: %s", peer))
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

		case leader := <-raftInstance.LeaderCh():
			if leader {
				debug.Log("raft", "I am the leader")
				go func() {
					time.Sleep(5 * time.Minute)
					debug.Log("raft", "Transferring leadership")
					raftInstance.LeadershipTransfer()
				}()
			} else {
				debug.Log("raft", "I am not the leader")
			}
		}
	}
}

func blockchainLoop(network *Network, raftInstance *raft.Raft, raftconsensus *libp2praft.Consensus, actor *libp2praft.Actor) {
	for {
		select {
		// New block added to the blockchain
		case <-raftconsensus.Subscribe():
			// First time the blockchain is updated
			newState, _ := raftconsensus.GetCurrentState()
			blockchain := newState.(*raftState).Blockchain
			debug.Log("raft", fmt.Sprintf("Blockchain updated, current length: %d", len(blockchain.Chain)))
			latestBlock := blockchain.GetLatestBlock()

			// Type assertion to access specific data
			switch latestBlock.BlockType {
			case "message":
				if messageData, ok := latestBlock.Data.(*models.MessageData); ok {
					debug.Log("raft", fmt.Sprintf("Latest message from: %s", messageData.Message.Sender))
					debug.Log("raft", fmt.Sprintf("Latest message: %s", messageData.Message.Message))
				}
			case "account":
				if accountData, ok := latestBlock.Data.(*models.AccountData); ok {
					debug.Log("raft", fmt.Sprintf("Latest account: %s", accountData.Account.Username))
				}
			case "firstMessage":
				if firstMessageData, ok := latestBlock.Data.(*models.FirstMessageData); ok {
					debug.Log("raft", fmt.Sprintf("Latest first message: %s and %s", firstMessageData.FirstMessage.PeerIDs[0], firstMessageData.FirstMessage.PeerIDs[1]))
				}
			default:
				debug.Log("raft", fmt.Sprintf("Latest block type: %s", latestBlock.BlockType))
			}
			network.ConsensusService.LatestBlock <- models.Block{
				Index:     latestBlock.Index,
				Timestamp: latestBlock.Timestamp,
				PrevHash:  latestBlock.PrevHash,
				Hash:      latestBlock.Hash,
				BlockType: latestBlock.BlockType,
				Data:      latestBlock.Data,
			}

		case inbound := <-network.PubSubService.Inbound:
			// If inbound is a message
			if message, ok := inbound.(models.Message); ok {
				debug.Log("raft", fmt.Sprintf("Inbound message: %s", message.Message))
				addMessageBlock(network, message, raftconsensus, actor)
			}
			if firstMessage, ok := inbound.(models.FirstMessage); ok {
				debug.Log("raft", fmt.Sprintf("Inbound first message: %s and %s", firstMessage.PeerIDs[0], firstMessage.PeerIDs[1]))
				addFirstMessageBlock(network, firstMessage, raftconsensus, actor)
			}
		}
	}
}

func addMessageBlock(network *Network, message models.Message, raftconsensus *libp2praft.Consensus, actor *libp2praft.Actor) {
	if actor.IsLeader() {
		debug.Log("raft", fmt.Sprintf("Adding message block: %s", message.Message))
		if message.Sender == "" || message.Receiver == "" || message.Message == "" {
			debug.Log("raft", "Message is missing required fields")
			return
		}
		if message.Sender == message.Receiver {
			debug.Log("raft", "Message sender and receiver cannot be the same")
			return
		}
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

func addFirstMessageBlock(network *Network, firstMessage models.FirstMessage, raftconsensus *libp2praft.Consensus, actor *libp2praft.Actor) {
	if actor.IsLeader() {
		sort.Strings(firstMessage.PeerIDs)
		if currentState := network.ConsensusService.Blockchain.Chain; currentState != nil {
			for _, block := range currentState {
				if block.BlockType == "firstMessage" {
					if block.Data.(*models.FirstMessageData).FirstMessage.PeerIDs[0] == firstMessage.PeerIDs[0] && block.Data.(*models.FirstMessageData).FirstMessage.PeerIDs[1] == firstMessage.PeerIDs[1] {
						debug.Log("raft", fmt.Sprintf("First message block already exists: %s and %s", block.Data.(*models.FirstMessageData).FirstMessage.PeerIDs[0], block.Data.(*models.FirstMessageData).FirstMessage.PeerIDs[1]))
						return
					}
				}
			}
		}
		if len(firstMessage.PeerIDs) != 2 {
			debug.Log("raft", "First message peer IDs must be exactly 2")
			return
		}
		if firstMessage.PeerIDs[0] == "" || firstMessage.PeerIDs[1] == "" {
			debug.Log("raft", "First message peer IDs cannot be empty")
			return
		}
		if firstMessage.SymetricKey0 == nil || firstMessage.SymetricKey1 == nil {
			debug.Log("raft", "First message symetric keys cannot be empty")
			return
		}

		debug.Log("raft", fmt.Sprintf("Adding first message block: %s and %s", firstMessage.PeerIDs[0], firstMessage.PeerIDs[1]))
		op := &raftOP{
			Type: "ADD_FIRST_MESSAGE_BLOCK",
			FirstMessage: &models.FirstMessage{
				PeerIDs:      firstMessage.PeerIDs,
				SymetricKey0: firstMessage.SymetricKey0,
				SymetricKey1: firstMessage.SymetricKey1,
			},
		}

		_, err := raftconsensus.CommitOp(op)
		if err != nil {
			debug.Log("err", fmt.Sprintf("Failed to commit block: %s", err))
		}
	}
}
