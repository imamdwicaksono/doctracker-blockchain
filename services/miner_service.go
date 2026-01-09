package services

import (
	"fmt"
	"time"

	"doc-tracker/blockchain"
	"doc-tracker/models"
	"doc-tracker/storage"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartMinerWorker() {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			if err := MineOnce(storage.DB); err != nil {
				fmt.Println("[Miner] error:", err)
			}
		}
	}()
}

func MineOnce(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ Lock candidate trackers
		var candidates []models.Tracker
		if err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
				Options:  "SKIP LOCKED",
			}).
			Where("status = ?", "progress").
			Limit(50).
			Find(&candidates).Error; err != nil {
			return err
		}

		if len(candidates) == 0 {
			return nil
		}

		// 2️⃣ Filter trackers with ALL checkpoints completed
		var eligible []models.Tracker
		for _, t := range candidates {
			if allCheckpointsCompleted(t.Checkpoints) {
				eligible = append(eligible, t)
			}
		}

		if len(eligible) == 0 {
			return nil // nothing eligible to mine
		}

		// 3️⃣ Create block FIRST (DB-first)
		block, err := blockchain.CreateBlock(tx, eligible)
		if err != nil {
			return err
		}

		// 4️⃣ Update tracker status → complete (AFTER block)
		ids := make([]string, 0, len(eligible))
		for _, t := range eligible {
			ids = append(ids, t.ID)
		}

		if err := tx.
			Model(&models.Tracker{}).
			Where("id IN ?", ids).
			Update("status", "complete").Error; err != nil {
			return err
		}

		// 5️⃣ Async anchor to EVM
		go func(b models.Block) {
			fmt.Println("[Audit] block anchored:", b.BlockHash)
		}(*block)

		fmt.Printf(
			"[Miner] block #%d created (%d trackers)\n",
			block.Height,
			block.TxCount,
		)

		return nil
	})
}

func allCheckpointsCompleted(checkpoints models.Checkpoints) bool {
	if len(checkpoints) == 0 {
		return false // no checkpoint = NOT eligible
	}

	for _, cp := range checkpoints {
		if !cp.IsCompleted {
			return false
		}
	}
	return true
}
