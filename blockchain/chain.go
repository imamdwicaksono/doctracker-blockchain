package blockchain

import (
	"doc-tracker/storage"
)

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
