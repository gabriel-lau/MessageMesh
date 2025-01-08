package models

type BlockType interface {
	GetType() string
}

type Block struct {
	BlockID           string     `json:"blockID"`
	CurrentBlockHash  string     `json:"currentBlockHash"`
	PreviousBlockHash string     `json:"previousBlockHash"`
	TimeStamp         string     `json:"timestamp"`
	BlockType         *BlockType `json:"blockType"`
}
