package blockchain

import (
	"crypto/rand"
	"crypto/sha256"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	Blockchain  []models.Block
	chainMutex  sync.RWMutex
	chainFile   = "data/chain.bin"  // File terenkripsi
	pubKeyPath  = "data/public.pem" // Sama dengan mempool
	privKeyPath = "data/private.pem"
)

// InitChain inisialisasi blockchain dengan genesis block terenkripsi
func InitChain() {
	chainMutex.Lock()
	defer chainMutex.Unlock()

	// Coba load dari storage dulu
	if loaded := loadChainFromStorage(); loaded {
		fmt.Printf("✅ Loaded blockchain with %d blocks\n", len(Blockchain))
		for _, block := range Blockchain {
			fmt.Printf("Block #%d | Hash: %s | Encrypted: %t | Transactions: %d\n", block.Index, block.Hash, block.Encrypted, len(block.Transactions))
			fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format(time.RFC3339))
			fmt.Printf("Previous Hash: %s\n", block.PrevHash)
			fmt.Println("Transactions:")
			for _, tx := range block.Transactions {
				fmt.Printf("  - ID: %s | Type: %s | Status: %s\n", tx.ID, tx.Type, tx.Status)
			}
			// Tambahkan garis pemisah antar block
			fmt.Println("--------------------------------------------------")
		}
		return
	}

	// Buat genesis block jika tidak ada
	genesis := CreateGenesisBlock()
	Blockchain = append(Blockchain, genesis)

	// Enkripsi dan simpan
	if err := saveEncryptedBlock(genesis); err != nil {
		log.Printf("⚠️ Failed to save genesis block: %v", err)
	}
}

// CreateGenesisBlock membuat block awal terenkripsi
func CreateGenesisBlock() models.Block {
	genesis := models.Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		PrevHash:     "0",
		Transactions: []models.Tracker{},
		Nonce:        0,
		Encrypted:    true,
	}
	genesis.Hash = CalculateHash(genesis)
	return genesis
}

// MineNewBlock membuat block baru terenkripsi
func MineNewBlock(transactions []models.Tracker) (models.Block, error) {
	prev := GetLastBlock()

	newBlock := models.Block{
		Index:        prev.Index + 1,
		Timestamp:    time.Now().Unix(),
		PrevHash:     prev.Hash,
		Transactions: transactions,
		Encrypted:    true,
	}

	// Mining process
	MineBlock(&newBlock, 4)

	// Enkripsi dan simpan block
	if err := saveEncryptedBlock(newBlock); err != nil {
		return models.Block{}, fmt.Errorf("failed to save block: %v", err)
	}

	chainMutex.Lock()
	Blockchain = append(Blockchain, newBlock)
	chainMutex.Unlock()

	for _, tx := range transactions {
		// Hapus tracker dari mempool
		mempool.RemoveFromMempool(tx.ID)
	}

	return newBlock, nil
}

// ================ ENCRYPTION FUNCTIONS ================

// saveEncryptedBlock menyimpan block terenkripsi
func saveEncryptedBlock(block models.Block) error {
	// 1. Serialisasi block
	data, err := json.Marshal(block)
	if err != nil {
		return fmt.Errorf("marshal failed: %v", err)
	}

	// 2. Load public key
	pubKey, err := utils.LoadECDSAPublicKey(pubKeyPath)
	if err != nil {
		return fmt.Errorf("load public key failed: %v", err)
	}

	// 3. Enkripsi dengan ECIES
	eciesPub := utils.ImportECDSAPublic(pubKey)
	encryptedData, err := utils.ECIESEncrypt(rand.Reader, eciesPub, data, nil, nil)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// 4. Simpan ke file terpisah per block
	blockFile := fmt.Sprintf("data/blocks/%d.bin", block.Index)
	if err := os.WriteFile(blockFile, encryptedData, 0600); err != nil {
		return fmt.Errorf("write failed: %v", err)
	}

	// 5. Update chain index
	return updateChainIndex(block.Index, block.Hash)
}

// loadDecryptedBlock memuat dan mendekripsi block
func loadDecryptedBlock(index int) (models.Block, error) {
	blockFile := fmt.Sprintf("data/blocks/%d.bin", index)
	encryptedData, err := os.ReadFile(blockFile)
	if err != nil {
		return models.Block{}, fmt.Errorf("read failed: %v", err)
	}

	// Load private key
	privKey, err := utils.LoadECDSAPrivateKey(privKeyPath)
	if err != nil {
		return models.Block{}, fmt.Errorf("load private key failed: %v", err)
	}

	// Dekripsi
	plaintext, err := utils.ECIESDecrypt(utils.ImportECDSA(privKey), encryptedData, nil, nil)
	if err != nil {
		return models.Block{}, fmt.Errorf("decryption failed: %v", err)
	}

	var block models.Block
	if err := json.Unmarshal(plaintext, &block); err != nil {
		return models.Block{}, fmt.Errorf("unmarshal failed: %v", err)
	}

	return block, nil
}

// updateChainIndex memperbarui index chain terenkripsi
func updateChainIndex(index int, hash string) error {
	chainData := fmt.Sprintf("%d:%s", index, hash)

	// Enkripsi data index
	pubKey, err := utils.LoadECDSAPublicKey(pubKeyPath)
	if err != nil {
		return err
	}

	eciesPub := utils.ImportECDSAPublic(pubKey)
	encryptedData, err := utils.ECIESEncrypt(rand.Reader, eciesPub, []byte(chainData), nil, nil)
	if err != nil {
		return err
	}

	// Tulis ke file chain index
	return os.WriteFile(chainFile, encryptedData, 0600)
}

// loadChainFromStorage memuat seluruh chain dari storage
func loadChainFromStorage() bool {
	if _, err := os.Stat(chainFile); os.IsNotExist(err) {
		return false
	}

	// Dekripsi chain index
	encryptedData, err := os.ReadFile(chainFile)
	if err != nil {
		log.Printf("⚠️ Failed to read chain file: %v", err)
		return false
	}

	privKey, err := utils.LoadECDSAPrivateKey(privKeyPath)
	if err != nil {
		log.Printf("⚠️ Failed to load private key: %v", err)
		return false
	}

	plaintext, err := utils.ECIESDecrypt(utils.ImportECDSA(privKey), encryptedData, nil, nil)
	if err != nil {
		log.Printf("⚠️ Failed to decrypt chain index: %v", err)
		return false
	}

	// Parse index terakhir
	plaintextStr, _ := utils.ConvertForAtoi(string(plaintext))
	if plaintextStr == "" {
		log.Println("⚠️ Invalid chain index format")
		return false
	}
	lastIndex, err := strconv.Atoi(plaintextStr)
	if err != nil {
		log.Printf("⚠️ Invalid chain index: %v", err)
		return false
	}

	// Load semua block
	var chain []models.Block
	for i := 0; i <= lastIndex; i++ {
		block, err := loadDecryptedBlock(i)
		if err != nil {
			log.Printf("⚠️ Failed to load block %d: %v", i, err)
			return false
		}
		chain = append(chain, block)
	}

	Blockchain = chain
	return true
}

// ================ CORE BLOCKCHAIN FUNCTIONS ================

// GetLastBlock mengambil block terakhir
func GetLastBlock() models.Block {
	chainMutex.RLock()
	defer chainMutex.RUnlock()

	if len(Blockchain) == 0 {
		return models.Block{}
	}
	return Blockchain[len(Blockchain)-1]
}

// IsBlockValid memvalidasi block
func IsBlockValid(newBlock, prevBlock models.Block) bool {
	if prevBlock.Index+1 != newBlock.Index {
		return false
	}
	if prevBlock.Hash != newBlock.PrevHash {
		return false
	}
	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}
	if !newBlock.Encrypted {
		return false
	}
	return true
}

// CalculateHash menghitung hash untuk block
func CalculateHash(block models.Block) string {
	record := strconv.Itoa(block.Index) + strconv.FormatInt(block.Timestamp, 10) + block.PrevHash + fmt.Sprintf("%v", block.Transactions) + strconv.Itoa(block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return fmt.Sprintf("%x", hashed)
}

// MineBlock melakukan proof-of-work
func MineBlock(block *models.Block, difficulty int) {
	for {
		hash := CalculateHash(*block)
		if hash[:difficulty] == "0000" {
			block.Hash = hash
			break
		}
		block.Nonce++
	}
}

// Iterate iterates over all transactions in the mempool and applies the given function.
func Iterate(fn func(tx *models.Tracker) error) error {
	for _, block := range GetAllBlocks() { // assuming mempool is a slice or map of *models.Transaction
		for i := range block.Transactions {
			if err := fn(&block.Transactions[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
