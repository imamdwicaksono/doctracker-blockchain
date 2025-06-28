package blockchain

import (
	"crypto/sha256"
	"doc-tracker/models"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

type Block = models.Block

func CalculateHash(block Block) string {
	data, _ := json.Marshal(struct {
		Index     int
		Timestamp int64
		PrevHash  string
		Nonce     int
		Tx        []models.Tracker
	}{
		Index:     block.Index,
		Timestamp: block.Timestamp,
		PrevHash:  block.PrevHash,
		Nonce:     block.Nonce,
		Tx:        block.Transactions,
	})
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

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

// MineBlock performs proof-of-work until hash matches the difficulty
func MineBlock(block *Block, difficulty int) {
	prefix := strings.Repeat("0", difficulty)
	for {
		block.Hash = CalculateHash(*block)
		if strings.HasPrefix(block.Hash, prefix) {
			break
		}
		block.Nonce++
	}
}
