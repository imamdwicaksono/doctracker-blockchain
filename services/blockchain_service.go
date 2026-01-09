package services

import (
	"doc-tracker/blockchain"
	"doc-tracker/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MinePendingTrackers(db *gorm.DB) (*models.Block, error) {

	var block *models.Block

	err := db.Transaction(func(tx *gorm.DB) error {

		// 1️⃣ Lock tracker rows
		var trackers []models.Tracker
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", "progress").
			Limit(50).
			Find(&trackers).Error; err != nil {
			return err
		}

		if len(trackers) == 0 {
			return nil // nothing to mine
		}

		// 2️⃣ Update tracker status
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

		// 3️⃣ Create audit block
		b, err := blockchain.CreateBlock(tx, trackers)
		if err != nil {
			return err
		}

		block = b
		return nil
	})

	return block, err
}
