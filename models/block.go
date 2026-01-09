package models

type Block struct {
	Height    int    `gorm:"primaryKey" json:"height"`
	BlockHash string `gorm:"size:66;uniqueIndex" json:"block_hash"`
	PrevHash  string `gorm:"size:66" json:"prev_hash"`
	TxCount   int    `json:"tx_count"`
	CreatedAt int64  `json:"created_at"`
}

type BlockchainAudit struct {
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	EvmTxHash   string `json:"evm_tx_hash"`
}
