package models

import (
	"MessageMesh/debug"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// BlockData interface defines common behavior for block data
type BlockData interface {
	CalculateDataHash() string
}

// Base Block struct
type Block struct {
	Index     int       `json:"Index"`
	Timestamp int64     `json:"Timestamp"`
	PrevHash  string    `json:"PrevHash"`
	Hash      string    `json:"Hash"`
	BlockType string    `json:"BlockType"`
	Data      BlockData `json:"Data"`
}

// MessageData implements BlockData
type MessageData struct {
	Message
}

func (md *MessageData) CalculateDataHash() string {
	return md.Sender + md.Receiver + md.Message.Message + md.Timestamp
}

type FirstMessageData struct {
	FirstMessage
}

func (md *FirstMessageData) CalculateDataHash() string {
	return md.PeerIDs[0] + md.PeerIDs[1] + hex.EncodeToString(md.SymetricKey0) + hex.EncodeToString(md.SymetricKey1)
}

// AccountData implements BlockData
type AccountData struct {
	Account
}

func (ad *AccountData) CalculateDataHash() string {
	return ad.Username + ad.PublicKey
}

// Updated CalculateHash method for Block
func (b *Block) CalculateHash() string {
	record := string(b.Index) + string(b.Timestamp) + b.PrevHash + b.BlockType
	if b.Data != nil {
		record += b.Data.CalculateDataHash()
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

type Blockchain struct {
	Chain []*Block
}

func CreateGenesisBlock() *Block {
	block := &Block{
		Index:     0,
		Timestamp: time.Now().Unix(),
		PrevHash:  "0",
		BlockType: "genesis",
		Data:      nil,
	}
	block.Hash = block.CalculateHash()
	return block
}

func (bc *Blockchain) GetMessageBlock(index int) *Block {
	block := bc.Chain[index]
	if block.BlockType != "message" {
		return nil
	}
	return block
}

func (bc *Blockchain) AddMessageBlock(message Message) *Block {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := &Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().Unix(),
		PrevHash:  prevBlock.Hash,
		BlockType: "message",
		Data:      &MessageData{Message: message},
	}
	newBlock.Hash = newBlock.CalculateHash()
	bc.Chain = append(bc.Chain, newBlock)
	return newBlock
}

func (bc *Blockchain) GetAccountBlock(index int) *Block {
	block := bc.Chain[index]
	if block.BlockType != "account" {
		return nil
	}
	return block
}

func (bc *Blockchain) AddAccountBlock(account Account) *Block {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := &Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().Unix(),
		PrevHash:  prevBlock.Hash,
		BlockType: "account",
		Data:      &AccountData{Account: account},
	}
	newBlock.Hash = newBlock.CalculateHash()
	bc.Chain = append(bc.Chain, newBlock)
	return newBlock
}

func (bc *Blockchain) AddFirstMessageBlock(firstMessage FirstMessage) *Block {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	sort.Strings(firstMessage.PeerIDs)
	newBlock := &Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().Unix(),
		PrevHash:  prevBlock.Hash,
		BlockType: "firstMessage",
		Data:      &FirstMessageData{FirstMessage: firstMessage},
	}
	newBlock.Hash = newBlock.CalculateHash()
	bc.Chain = append(bc.Chain, newBlock)
	return newBlock
}

func (bc *Blockchain) GetFirstMessageBlock(index int) *Block {
	block := bc.Chain[index]
	if block.BlockType != "firstMessage" {
		return nil
	}
	return block
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		if currentBlock.Hash != currentBlock.CalculateHash() {
			return false
		}

		if currentBlock.PrevHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

// Check if the blockchain has a first message block with a specific peer
func (bc *Blockchain) CheckPeerFirstMessage(peerIDs []string) *FirstMessage {
	// Loop through the blockchain
	for _, block := range bc.Chain {
		if block.BlockType == "firstMessage" {
			sort.Strings(peerIDs)
			if block.Data.(*FirstMessageData).PeerIDs[0] == peerIDs[0] && block.Data.(*FirstMessageData).PeerIDs[1] == peerIDs[1] {
				// Check if the first message is expired
				if block.Data.(*FirstMessageData).Expiry < time.Now().Unix() {
					debug.Log("blockchain", fmt.Sprintf("First message expired for %s and %s", peerIDs[0], peerIDs[1]))
					continue
				}
				debug.Log("blockchain", fmt.Sprintf("First message found for %s and %s", peerIDs[0], peerIDs[1]))
				return &block.Data.(*FirstMessageData).FirstMessage
			}
		}
	}
	return nil
}
