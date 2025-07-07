package blockchain

import (
	"doc-tracker/models"
	"time"
)

type Block = models.Block

// NewBlock creates a new block from completed transactions
func NewBlock(index int, prevHash string, trackers []models.Tracker, timestamp int64, nonce int) Block {
	block := Block{
		Index:        index,
		Timestamp:    timestamp,
		PrevHash:     prevHash,
		Nonce:        nonce,
		Transactions: trackers,
	}
	block.Hash = CalculateHash(block)
	return block
}

func NewBlockFromTransactions(trackers []models.Tracker) Block {
	lastBlock := GetLastBlock()
	index := lastBlock.Index + 1
	timestamp := time.Now().Unix()
	nonce := 0 // jika nanti perlu Proof of Work
	block := Block{
		Index:        index,
		Timestamp:    timestamp,
		PrevHash:     lastBlock.Hash,
		Nonce:        nonce,
		Transactions: trackers,
	}
	block.Hash = CalculateHash(block)
	return block
}
