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
	// Start progress
	a.network.Progress = backend.NewProgress()

	a.ctx = ctx

	// Start the network, connect to peers and join the blockchain
	a.network.ConnectToNetwork()

	// Start the UI loop
	go backend.UIDataLoop(a.network, a.ctx)
}

// CHATCOMPONET
func (a *App) SendMessage(message string, receiver string) {
	a.network.SendMessage(message, receiver)
}

func (a *App) GetBlockchain() []*models.Block {
	return a.network.ConsensusService.Blockchain.Chain
}

func (a *App) GetMessages() []*models.Message {
	messages := make([]*models.Message, 0)
	for _, block := range a.network.ConsensusService.Blockchain.Chain {
		if block.BlockType == "message" {
			messages = append(messages, &block.Data.(*models.MessageData).Message)
		}
	}
	return messages
}

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
func (a *App) GetAccounts() []*models.Account {
	accounts := make([]*models.Account, 0)
	for _, block := range a.network.ConsensusService.Blockchain.Chain {
		if block.BlockType == "account" {
			accounts = append(accounts, &block.Data.(*models.AccountData).Account)
		}
	}
	return accounts
}
