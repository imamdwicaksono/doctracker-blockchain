package blockchain

import (
	"doc-tracker/models"
	"doc-tracker/storage"
	"time"
)

var Blockchain []models.Block

// CreateGenesisBlock membuat block awal
func CreateGenesisBlock() models.Block {
	genesis := models.Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		PrevHash:     "0",
		Transactions: []models.Tracker{},
		Nonce:        0,
	}
	genesis.Hash = CalculateHash(genesis)
	return genesis
}

// InitChain inisialisasi blockchain dengan genesis block
func InitChain() {
	genesis := CreateGenesisBlock()
	Blockchain = append(Blockchain, genesis)
	storage.SaveBlock(genesis)
}

// MineNewBlock membuat block baru dan memproses mining
func MineNewBlock(transactions []models.Tracker) Block {
	prev := GetLastBlock()

	newBlock := Block{
		Index:        prev.Index + 1,
		Timestamp:    time.Now().Unix(),
		PrevHash:     prev.Hash,
		Transactions: transactions,
	}

	// Modular mining
	MineBlock(&newBlock, 4)

	Blockchain = append(Blockchain, newBlock)
	storage.SaveBlock(newBlock)

	return newBlock
}

// AddBlock menambahkan block dari peer lain jika belum ada
func AddBlock(block Block) {
	for _, b := range Blockchain {
		if b.Hash == block.Hash {
			return // sudah ada
		}
	}
	Blockchain = append(Blockchain, block)
	storage.SaveBlock(block)
}

// GetLastBlock mengambil block terakhir
func GetLastBlock() Block {
	return Blockchain[len(Blockchain)-1]
}

// IsBlockValid memvalidasi block baru dari peer
func IsBlockValid(newBlock, prevBlock Block) bool {
	if prevBlock.Index+1 != newBlock.Index {
		return false
	}
	if prevBlock.Hash != newBlock.PrevHash {
		return false
	}
	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}
	if newBlock.Hash[:4] != "0000" {
		return false
	}
	return true
}

// IsValidChain memeriksa apakah chain valid dari genesis hingga terakhir
func IsValidChain(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		if !IsBlockValid(chain[i], chain[i-1]) {
			return false
		}
	}
	return true
}

// ReplaceChain mengganti blockchain lokal jika newChain lebih panjang dan valid
func ReplaceChain(newChain []Block) bool {
	if len(newChain) <= len(Blockchain) {
		return false
	}

	if !IsValidChain(newChain) {
		return false
	}

	Blockchain = newChain

	// Simpan semua block baru
	storage.ClearAllBlocks()
	for _, b := range newChain {
		storage.SaveBlock(b)
	}

	return true
}

// GetAllBlocks mengembalikan seluruh blockchain
func GetAllBlocks() []Block {
	return Blockchain
}

func LoadChainFromStorage() {
	Blockchain = storage.LoadAllBlocks()

	if len(Blockchain) == 0 {
		InitChain() // kalau belum ada sama sekali, buat Genesis Block
	}
}

// TryAddBlock menambahkan block jika belum ada dan valid
func TryAddBlock(block Block) bool {
	last := GetLastBlock()

	if block.Hash == last.Hash || block.Index <= last.Index {
		return false // sudah ada atau lebih lama
	}

	if IsBlockValid(block, last) {
		Blockchain = append(Blockchain, block)
		storage.SaveBlock(block)
		return true
	}

	return false // tidak valid
}

func AddBlockToChain(block Block) {
	last := GetLastBlock()
	if IsBlockValid(block, last) {
		Blockchain = append(Blockchain, block)
		storage.SaveBlock(block)
	}
}
