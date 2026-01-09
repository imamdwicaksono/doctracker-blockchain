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
			if err := mineOnce(storage.DB); err != nil {
				fmt.Println("[Miner] error:", err)
			}
		}
	}()
}

func mineOnce(db *gorm.DB) error {

	return db.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ Lock tracker rows (multi-worker safe)
		var trackers []models.Tracker
		if err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
				Options:  "SKIP LOCKED",
			}).
			Where("status = ?", "progress").
			Limit(50).
			Find(&trackers).Error; err != nil {
			return err
		}

		if len(trackers) == 0 {
			return nil // nothing to mine
		}

		// 2️⃣ Update tracker status → complete
		ids := make([]string, 0, len(trackers))
		for _, t := range trackers {
			ids = append(ids, t.ID)
		}

		if err := tx.
			Model(&models.Tracker{}).
			Where("id IN ?", ids).
			Update("status", "complete").Error; err != nil {
			return err
		}

		// 3️⃣ Create audit block (DB-first)
		block, err := blockchain.CreateBlock(tx, trackers)
		if err != nil {
			return err
		}

		// 4️⃣ Async EVM audit (non-blocking)
		go func(b models.Block) {
			// TODO: kirim hash ke EVM & SaveAudit
			fmt.Println("[Audit] block anchored:", b.BlockHash)
		}(*block)

		fmt.Printf("[Miner] block #%d created (%d trackers)\n",
			block.Height, block.TxCount)

		return nil
	})
}
