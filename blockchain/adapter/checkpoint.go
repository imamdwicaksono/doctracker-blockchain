package adapter

import (
	"gorm.io/gorm"
)

func LoadLastBlock(db *gorm.DB) uint64 {
	var block uint64
	db.Raw(
		"SELECT last_block FROM blockchain_checkpoint WHERE id = 1",
	).Scan(&block)
	return block
}

func SaveLastBlock(db *gorm.DB, block uint64) {
	_ = db.Exec(
		"UPDATE blockchain_checkpoint SET last_block = $1 WHERE id = 1",
		block,
	)
}
