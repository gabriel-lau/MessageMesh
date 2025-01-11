package backend

import (
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	consensus "github.com/libp2p/go-libp2p-consensus"
	libp2praft "github.com/libp2p/go-libp2p-raft"
)

type raftState struct {
	Now string
}

type raftOP struct {
	Op string
}

func (o *raftOP) ApplyTo(state consensus.State) (consensus.State, error) {
	fmt.Println("Applying OP: ", o.Op)
	return state, nil
}

func StartRaft(network *Network) {
	pids := network.ChatRoom.PeerList()
	// pids := network.P2p.Host.Peerstore().Peers()
	// -- Create the consensus with no actor attached
	raftconsensus := libp2praft.NewOpLog(&raftState{}, &raftOP{})
	// raftconsensus = libp2praft.NewConsensus(&raftState{"i am not consensuated"})
	// --

	// -- Create Raft servers configuration
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
	// --

	// -- Create LibP2P transports Raft
	transport, err := libp2praft.NewLibp2pTransport(network.P2p.Host, 2*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	// --

	// -- Configuration
	raftQuiet := true
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
		fmt.Println(err)
	}

	// -- Log store and stable store: we use inmem.
	// logStore := raft.NewInmemStore()
	logStore, _ := raftboltdb.NewBoltStore("db/raft.db")
	// --

	// -- Boostrap everything if necessary
	bootstrapped, err := raft.HasExistingState(logStore, logStore, snapshots)
	if err != nil {
		fmt.Println(err)
	}

	if !bootstrapped {
		// Bootstrap cluster first
		raft.BootstrapCluster(config, logStore, logStore, snapshots, transport, serverConfig)
	} else {
		fmt.Println("Already initialized!!")
	}
	// --

	// Create Raft instance. Our consensus.FSM() provides raft.FSM
	// implementation
	raft, err := raft.NewRaft(config, raftconsensus.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		fmt.Println(err)
	}

	actor := libp2praft.NewActor(raft)
	raftconsensus.SetActor(actor)

	go func() {
		for {
			if actor.IsLeader() {
				fmt.Println("I am the leader")
				fmt.Println("Raft State: " + raft.State().String())
				updateState(raftconsensus)
			} else {
				fmt.Println("I am not the leader")
				fmt.Println("Raft State: " + raft.State().String())
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func updateState(c *libp2praft.Consensus) {
	loc, _ := time.LoadLocation("UTC")
	newState := &raftState{Now: time.Now().In(loc).String()}

	// CommitState() blocks until the state has been
	// agreed upon by everyone
	agreedState, err := c.CommitState(newState)
	if err != nil {
		fmt.Println(err)
	}
	if agreedState == nil {
		fmt.Println("agreedState is nil: commited on a non-leader?")
	}
	agreedRaftState := agreedState.(*raftState)

	fmt.Printf("Current state: %d\n", agreedRaftState.Now)
}
