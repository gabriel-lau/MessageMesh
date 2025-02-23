package models

import (
	"crypto/sha256"
	"encoding/hex"
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
	return md.SymetricKey
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
