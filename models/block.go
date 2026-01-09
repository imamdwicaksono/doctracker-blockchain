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
