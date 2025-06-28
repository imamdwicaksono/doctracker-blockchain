package storage

import (
	"doc-tracker/models"
	"encoding/json"
	"os"
)

const ChainFile = "data/blockchain.json"

type Block = models.Block

func SaveBlock(block Block) {
	blocks := LoadAllBlocks()
	blocks = append(blocks, block)

	file, _ := os.Create(ChainFile)
	defer file.Close()
	_ = json.NewEncoder(file).Encode(blocks)
}

func LoadAllBlocks() []Block {
	var blocks []Block
	file, err := os.Open(ChainFile)
	if err != nil {
		return []Block{}
	}
	defer file.Close()

	_ = json.NewDecoder(file).Decode(&blocks)
	return blocks
}
