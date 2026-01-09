package models

type Block struct {
	Height    int64  `gorm:"column:height;primaryKey;autoIncrement:false"`
	BlockHash string `gorm:"size:64;not null;uniqueIndex"`
	PrevHash  string `gorm:"size:64"`
	TxCount   int    `gorm:"not null"`
	CreatedAt int64  `gorm:"not null"`
}

type BlockchainAudit struct {
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	EvmTxHash   string `json:"evm_tx_hash"`
}

type BlockPayload struct {
	Height     uint64   `json:"height"`
	PrevHash   string   `json:"prev_hash"`
	Timestamp  int64    `json:"timestamp"`
	TxHashes   []string `json:"tx_hashes"`
	MerkleRoot string   `json:"merkle_root"`
	BlockHash  string   `json:"block_hash"`
}
