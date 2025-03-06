package backend

import (
	"MessageMesh/backend/models"
	"MessageMesh/debug"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func UIDataLoop(network Network, ctx context.Context) {
	debug.Log("ui", "Wails events emitter started")
	// Emit ready event
	if !debug.IsHeadless {
		runtime.EventsEmit(ctx, "ready")
		// Send the user's peer ID once to the frontend and then remove the event listener
		runtime.EventsEmit(ctx, "getUserPeerID", network.P2pService.Host.ID())
		runtime.EventsEmit(ctx, "getPeerList", network.PubSubService.PeerList())
		for {
			select {
			case peerIDs := <-network.PubSubService.PeerIDs:
				runtime.EventsEmit(ctx, "getPeerList", peerIDs)
				debug.Log("ui", "Peers: "+string(len(peerIDs)))

			case connected := <-network.ConsensusService.Connected:
				runtime.EventsEmit(ctx, "getConnected", connected)
				debug.Log("ui", "Consensus connected: "+fmt.Sprint(connected))

			// repeat this every 10 seconds
			case <-time.After(10 * time.Second):
				runtime.EventsEmit(ctx, "getPeerList", network.PubSubService.PeerList())
				runtime.EventsEmit(ctx, "getConnected", network.ConsensusService.Connected)

			case block := <-network.ConsensusService.LatestBlock:
				// Check if the block is a message block
				if block.BlockType == "message" {
					peerID := network.PubSubService.selfid.String()
					// If the message is encrypted, decrypt it
					if block.Data.(*models.MessageData).Sender == peerID || block.Data.(*models.MessageData).Receiver == peerID {
						decryptedMessage, err := network.DecryptMessage(block.Data.(*models.MessageData).Message.Message, block.Data.(*models.MessageData).Message.Sender)
						if err != nil {
							debug.Log("ui", "Error decrypting message: "+err.Error())
						}
						block.Data.(*models.MessageData).Message.Message = decryptedMessage
					}
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

				// Get the blockchain
				runtime.EventsEmit(ctx, "getBlockchain", network.ConsensusService.Blockchain.Chain)

				// Get the messages
				messages := make([]*models.Message, 0)
				for _, block := range network.ConsensusService.Blockchain.Chain {
					if block.BlockType == "message" {
						messages = append(messages, &block.Data.(*models.MessageData).Message)
					}
				}
				debug.Log("ui", "Messages: "+string(len(messages)))
				runtime.EventsEmit(ctx, "getMessages", messages)
				// Get the accounts
				accounts := make([]*models.Account, 0)
				for _, block := range network.ConsensusService.Blockchain.Chain {
					if block.BlockType == "account" {
						accounts = append(accounts, &block.Data.(*models.AccountData).Account)
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
				network.SendEncryptedMessage("Hello I am "+debug.Username, "Qma9HU4gynWXNzWwpqmHRnLXikstTgCbYHfG6aqJTLrxfq")
			case <-ctx.Done():
				return
			}
		}
	}
}
