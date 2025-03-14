package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"context"
	"encoding/hex"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func UIDataLoop(network Network, ctx context.Context) {
	debug.Log("ui", "Wails events emitter started")
	if !debug.IsHeadless {
		runtime.EventsEmit(ctx, "getUserPeerID", network.P2pService.Host.ID())
		runtime.EventsEmit(ctx, "getPeerList", network.PubSubService.PeerList())
		for {
			select {
			case peerIDs := <-network.PubSubService.PeerIDs:
				runtime.EventsEmit(ctx, "getPeerList", peerIDs)
				debug.Log("ui", "Peers: "+string(len(peerIDs)))
				runtime.EventsEmit(ctx, "getConnected", true)

			// repeat this every 10 seconds
			case <-time.After(10 * time.Second):
				runtime.EventsEmit(ctx, "getPeerList", network.PubSubService.PeerList())

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
				if block.BlockType == "firstMessage" {
					runtime.EventsEmit(ctx, "getFirstMessage", block.Data.(*models.FirstMessageData).FirstMessage)
					debug.Log("ui", "First Message: "+hex.EncodeToString(block.Data.(*models.FirstMessageData).FirstMessage.SymetricKey0)+" and "+hex.EncodeToString(block.Data.(*models.FirstMessageData).FirstMessage.SymetricKey1))
				}

				runtime.EventsEmit(ctx, "getBlock", block)
				runtime.EventsEmit(ctx, "getBlockchain", network.ConsensusService.Blockchain.Chain)

			case <-network.ConsensusService.Connected:
				runtime.EventsEmit(ctx, "getConnected", true)
			}
		}
	} else {
		for {
			select {
			case block := <-network.ConsensusService.LatestBlock:
				debug.Log("ui", "Block: "+block.BlockType)
			// case <-time.After(30 * time.Second):
			// 	network.SendEncryptedMessage("Its "+time.Now().Format("2006-01-02 15:04:05")+" I am "+debug.Username, "Qma9HU4gynWXNzWwpqmHRnLXikstTgCbYHfG6aqJTLrxfq")
			case <-ctx.Done():
				return
			}
		}
	}
}
