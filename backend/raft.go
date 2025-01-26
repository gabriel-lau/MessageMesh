package backend

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/raft"
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
	pids = append(pids, network.P2p.Host.ID())
	servers := make([]raft.Server, len(pids))
	for i, pid := range pids {
		servers[i] = raft.Server{
			Suffrage: raft.Voter,
			ID:       raft.ServerID(pid.String()),
			Address:  raft.ServerAddress(pid.String()),
		}
		fmt.Println("Server: ", servers[i])
	}
	serverConfig := raft.Configuration{
		Servers: servers,
	}
	// --

	// -- Create LibP2P transports Raft
	transport, err := libp2praft.NewLibp2pTransport(network.P2p.Host, 3*time.Second)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}

	// -- Log store and stable store: we use inmem.
	logStore := raft.NewInmemStore()
	// logStore, _ := raftboltdb.NewBoltStore("db/raft.db")
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

	raftInstance, err := raft.NewRaft(config, raftconsensus.FSM(), logStore, logStore, snapshots, transport)
	if err != nil {
		fmt.Println(err)
	}

	actor := libp2praft.NewActor(raftInstance)
	raftconsensus.SetActor(actor)

	waitForLeader(raftInstance)

	go func() {
		refreshticker := time.NewTicker(time.Second)
		defer refreshticker.Stop()
		for {
			select {
			case <-raftconsensus.Subscribe():
				newState, _ := raftconsensus.GetCurrentState()
				fmt.Println("New state is: ", newState.(*raftState).Now)

			case <-raftInstance.LeaderCh():
				fmt.Println("Leader changed")

			case <-refreshticker.C:
				fmt.Println("Number of peers in network: ", network.ChatRoom.PeerList())
				updateConnectedServers(network, raftInstance, servers)

				if actor.IsLeader() {
					fmt.Println("I am the leader")
					fmt.Println("Raft State: " + raftInstance.State().String())
					fmt.Println(("Number of peers: "), raftInstance.Stats()["num_peers"])
					updateState(raftconsensus)
					getState(raftconsensus)
				} else {
					fmt.Println("I am not the leader")
					fmt.Println("Leader is: ", raftInstance.Leader())
					fmt.Println("Raft State: " + raftInstance.State().String())
					fmt.Println(("Number of peers: "), raftInstance.Stats()["num_peers"])
					getState(raftconsensus)
				}
			}
		}
	}()

	go updateConnectedServers(network, raftInstance, servers)
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
}

func getState(c *libp2praft.Consensus) {
	state, err := c.GetCurrentState()
	if err != nil {
		fmt.Println(err)
	}
	if state == nil {
		fmt.Println("state is nil: commited on a non-leader?")
		return
	}
	fmt.Printf("Current state: %d\n", state)
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
			fmt.Println("timed out waiting for Leader")
			fmt.Println("Current Raft State: ", r.State())
			fmt.Println("Current Leader: ", r.Leader())
			return
		}
	}
}

func updateConnectedServers(network *Network, raftInstance *raft.Raft, servers []raft.Server) {
	// Listen for peer joins and leaves
	for {
		select {
		case peer := <-network.ChatRoom.PeerJoin:
			fmt.Println("Peer joined: ", peer)
			raftInstance.AddVoter(raft.ServerID(peer.String()), raft.ServerAddress(peer.String()), 0, 0)
		case peer := <-network.ChatRoom.PeerLeave:
			fmt.Println("Peer left: ", peer)
			raftInstance.RemoveServer(raft.ServerID(peer.String()), 0, 0)
		}
	}
}
