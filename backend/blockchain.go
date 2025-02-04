package backend

import (
	"MessageMesh/backend/models"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Base Block struct
type Block struct {
	Index     int
	Timestamp int64
	PrevHash  string
	Hash      string
	BlockType string // "message" or "account"
}

// MessageBlock extends Block
type MessageBlock struct {
	Block
	Message *models.Message
}

// AccountBlock extends Block
type AccountBlock struct {
	Block
	Account *models.Account
}

type Blockchain struct {
	Chain []Block
}

func (b *Block) CalculateHash() string {
	record := string(b.Index) + string(b.Timestamp) + b.PrevHash + b.BlockType
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (mb *MessageBlock) CalculateHash() string {
	record := string(mb.Index) + string(mb.Timestamp) + mb.PrevHash +
		mb.Message.Sender + mb.Message.Receiver + mb.Message.Message + mb.Message.Timestamp
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (ab *AccountBlock) CalculateHash() string {
	record := string(ab.Index) + string(ab.Timestamp) + ab.PrevHash +
		ab.Account.Username + ab.Account.PublicKey
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CreateGenesisBlock() Block {
	block := Block{
		Index:     0,
		Timestamp: time.Now().Unix(),
		PrevHash:  "0",
		BlockType: "genesis",
	}
	block.Hash = block.CalculateHash()
	return block
}

func (bc *Blockchain) GetMessageBlock(index int) *MessageBlock {
	block := bc.Chain[index]
	if block.BlockType != "message" {
		return nil
	}
	messageBlock := &MessageBlock{}
	messageBlock.Block = block
	return messageBlock
}

func (bc *Blockchain) AddMessageBlock(message *models.Message) *MessageBlock {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := &MessageBlock{
		Block: Block{
			Index:     prevBlock.Index + 1,
			Timestamp: time.Now().Unix(),
			PrevHash:  prevBlock.Hash,
			BlockType: "message",
		},
		Message: message,
	}
	newBlock.Hash = newBlock.CalculateHash()
	bc.Chain = append(bc.Chain, newBlock.Block)
	return newBlock
}

func (bc *Blockchain) GetAccountBlock(index int) *AccountBlock {
	block := bc.Chain[index]
	if block.BlockType != "account" {
		return nil
	}
	accountBlock := &AccountBlock{}
	accountBlock.Block = block
	return accountBlock
}

func (bc *Blockchain) AddAccountBlock(account *models.Account) *AccountBlock {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := &AccountBlock{
		Block: Block{
			Index:     prevBlock.Index + 1,
			Timestamp: time.Now().Unix(),
			PrevHash:  prevBlock.Hash,
			BlockType: "account",
		},
		Account: account,
	}
	newBlock.Hash = newBlock.CalculateHash()
	bc.Chain = append(bc.Chain, newBlock.Block)
	return newBlock
}

func (bc *Blockchain) GetLatestBlock() Block {
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
