package adapter

import "database/sql"

func LoadLastBlock(db *sql.DB) uint64 {
	var block uint64
	_ = db.QueryRow(
		"SELECT last_block FROM blockchain_checkpoint WHERE id = 1",
	).Scan(&block)
	return block
}

func SaveLastBlock(db *sql.DB, block uint64) {
	_, _ = db.Exec(
		"UPDATE blockchain_checkpoint SET last_block = $1 WHERE id = 1",
		block,
	)
}
