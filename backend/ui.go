package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func UIDataLoop(network Network, ctx context.Context) {
	debug.Log("ui", "Wails events emitter started")
	if !debug.IsHeadless {

		// Send the user's peer ID once to the frontend and then remove the event listener
		runtime.EventsEmit(ctx, "getUserPeerID", network.P2pService.Host.ID())
		runtime.EventsEmit(ctx, "getPeerList", network.PubSubService.PeerList())
		for {
			select {
			case peerIDs := <-network.PubSubService.PeerIDs:
				runtime.EventsEmit(ctx, "getPeerList", peerIDs)
				debug.Log("ui", "Peers: "+string(len(peerIDs)))
			case block := <-network.ConsensusService.LatestBlock:
				// Check if the block is a message block
				if block.BlockType == "message" {
					runtime.EventsEmit(ctx, "getMessage", block.Data.(*models.MessageData).Message)
					debug.Log("ui", "Message: "+block.Data.(*models.MessageData).Message.Message)
				}
				if block.BlockType == "account" {
					runtime.EventsEmit(ctx, "getAccount", block.Data.(*models.AccountData).Account)
					debug.Log("ui", "Account: "+block.Data.(*models.AccountData).Account.Username)
				}
				// Get the blockchain
				runtime.EventsEmit(ctx, "getBlockchain", network.ConsensusService.Blockchain.Chain)

				// Get the messages
				messages := make([]*models.Message, 0)
				for _, block := range network.ConsensusService.Blockchain.Chain {
					if block.BlockType == "message" {
						messages = append(messages, block.Data.(*models.MessageData).Message)
					}
				}
				runtime.EventsEmit(ctx, "getMessages", messages)
				// Get the accounts
				accounts := make([]*models.Account, 0)
				for _, block := range network.ConsensusService.Blockchain.Chain {
					if block.BlockType == "account" {
						accounts = append(accounts, block.Data.(*models.AccountData).Account)
					}
				}
				runtime.EventsEmit(ctx, "getAccounts", accounts)
			}
		}
	} else {
		for {
			select {
			case block := <-network.ConsensusService.LatestBlock:
				debug.Log("ui", "Block: "+block.BlockType)
			case <-time.After(30 * time.Second):
				network.SendMessage("Hello I am "+debug.Username, "Qma9HU4gynWXNzWwpqmHRnLXikstTgCbYHfG6aqJTLrxfq")
			case <-ctx.Done():
				return
			}
		}
	}
}
