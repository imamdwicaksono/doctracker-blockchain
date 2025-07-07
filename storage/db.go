package storage

import (
	"doc-tracker/models"
	"doc-tracker/utils"
	"encoding/json"
	"os"
)

const (
	plainFile  = "data/blockchain.json"
	encFile    = "data/blockchain.json.enc"
	privateKey = "data/private.pem"
	publicKey  = "data/public.pem"
)

type Block = models.Block

func LoadAllBlocks() []Block {
	var blocks []Block

	// Decrypt file terenkripsi ke file sementara
	err := utils.DecryptFileWithPEM(privateKey, encFile, plainFile)
	if err != nil {
		return []Block{} // fallback: kosong jika gagal decrypt
	}
	defer os.Remove(plainFile) // aman, kita hapus setelah baca

	// Baca dari file hasil decrypt
	file, err := os.Open(plainFile)
	if err != nil {
		return []Block{}
	}
	defer file.Close()

	_ = json.NewDecoder(file).Decode(&blocks)
	return blocks
}
