package main

import (
	backend "MessageMesh/backend"
	"MessageMesh/backend/models"
	"context"
)

// App struct
type App struct {
	ctx     context.Context
	network backend.Network
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start the network, connect to peers and join the blockchain
	a.network.ConnectToNetwork()

	// Start the UI Data loop
	go backend.UIDataLoop(a.network, a.ctx)
}

// Functions for the UI to get data from the network

// Get the list of peers in the network
func (a *App) GetPeerList() []string {
	peers := make([]string, 0)
	for _, peer := range a.network.PubSubService.PeerList() {
		peers = append(peers, peer.String())
	}
	return peers
}

// Get the user's peer ID
func (a *App) GetUserPeerID() string {
	return a.network.PubSubService.SelfID().String()
}

// Send a message to a peer (Not in use by the UI)
func (a *App) SendMessage(message string, receiver string) {
	a.network.SendMessage(message, receiver)
}

// Send an encrypted message to a peer
func (a *App) SendEncryptedMessage(message string, receiver string) {
	a.network.SendEncryptedMessage(message, receiver)
}

// Get the blockchain (Not in use by the UI)
func (a *App) GetBlockchain() []*models.Block {
	return a.network.ConsensusService.Blockchain.Chain
}

// Get the messages from the blockchain (Not in use by the UI)
func (a *App) GetMessages() []*models.Message {
	messages := make([]*models.Message, 0)
	for _, block := range a.network.ConsensusService.Blockchain.Chain {
		if block.BlockType == "message" {
			messages = append(messages, &block.Data.(*models.MessageData).Message)
		}
	}
	return messages
}

// Get a decrypted message from the blockchain
func (a *App) GetDecryptedMessage(message string, peerIDs []string) (string, error) {
	return a.network.DecryptMessage(message, peerIDs)
}

// Get the messages from a specific peer
func (a *App) GetMessagesFromPeer(peer string) []*models.Message {
	messages := make([]*models.Message, 0)
	for _, block := range a.network.ConsensusService.Blockchain.Chain {
		if block.BlockType == "message" {
			if block.Data.(*models.MessageData).Message.Sender == peer || block.Data.(*models.MessageData).Message.Receiver == peer {
				messages = append(messages, &block.Data.(*models.MessageData).Message)
			}
		}
	}
	return messages
}

// Get the accounts from the blockchain (Not in use by the UI)
func (a *App) GetAccounts() []*models.Account {
	accounts := make([]*models.Account, 0)
	for _, block := range a.network.ConsensusService.Blockchain.Chain {
		if block.BlockType == "account" {
			accounts = append(accounts, &block.Data.(*models.AccountData).Account)
		}
	}
	return accounts
}

func (a *App) SetTopic(topic string) {
	//a.network.PubSubService.Topic = topic
}
