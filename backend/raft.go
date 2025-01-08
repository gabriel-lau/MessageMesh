package backend

import (
	"MessageMesh/backend/models"

	libp2praft "github.com/libp2p/go-libp2p-raft"
)

func StartRaft(host *P2P) {
	consensus := libp2praft.NewConsensus(&models.Block{"0", "0", "0", "0", nil})
	consensus.FSM()

}
