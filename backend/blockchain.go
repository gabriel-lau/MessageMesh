package backend

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index        int
	Timestamp    int64
	Data         string
	PrevHash     string
	Hash         string
	Transactions []Transaction
}

type Transaction struct {
	From   string
	To     string
	Amount float64
}

type Blockchain struct {
	Chain []Block
}

func (b *Block) CalculateHash() string {
	record := string(b.Index) + string(b.Timestamp) + b.Data + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CreateGenesisBlock() Block {
	return Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Data:         "Genesis Block",
		PrevHash:     "0",
		Transactions: []Transaction{},
	}
}

func (bc *Blockchain) AddBlock(data string, transactions []Transaction) Block {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PrevHash:     prevBlock.Hash,
		Transactions: transactions,
	}
	newBlock.Hash = newBlock.CalculateHash()
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
