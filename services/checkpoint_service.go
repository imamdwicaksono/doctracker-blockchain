package services

import (
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/storage"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateCheckpointStatus(
	trackerID string,
	checkpointAddr string,
	evidenceHash string,
	evidencePath string,
) error {

	return storage.DB.Transaction(func(tx *gorm.DB) error {

		var tracker models.Tracker
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", trackerID).
			First(&tracker).Error; err != nil {
			return err
		}

		updated := false

		for i := range tracker.Checkpoints {
			cp := &tracker.Checkpoints[i]

			for _, addr := range cp.Addresses {
				if addr == checkpointAddr {

					if cp.IsCompleted {
						return fmt.Errorf("checkpoint already completed")
					}

					cp.IsCompleted = true
					cp.CompletedAt = time.Now().Unix()
					cp.EvidenceHash = evidenceHash
					cp.EvidencePath = evidencePath
					cp.Note = ""

					updated = true
					break
				}
			}
		}

		if !updated {
			return fmt.Errorf("checkpoint not found")
		}

		// 🔥 overwrite JSONB explicitly
		if err := tx.
			Model(&models.Tracker{}).
			Where("id = ?", trackerID).
			Update("checkpoints", tracker.Checkpoints).
			Error; err != nil {
			return err
		}

		// 🔁 cek apakah semua checkpoint complete
		allComplete := true
		for _, cp := range tracker.Checkpoints {
			if !cp.IsCompleted {
				allComplete = false
				break
			}
		}

		newStatus := "progress"
		if allComplete {
			newStatus = "complete"
		}

		if err := tx.
			Model(&models.Tracker{}).
			Where("id = ?", trackerID).
			Update("status", newStatus).
			Error; err != nil {
			return err
		}

		// 🔄 sync mempool cache
		mempool.UpdateTracker(tx, &tracker)

		// 🚀 trigger miner jika complete
		if allComplete {
			go StartMinerWorker()
		}

		return nil
	})
}

func GetCheckpointAddressByEmail(trackerID, email string) string {

	tracker, err := GetTrackerByID(trackerID)
	if err != nil {
		return ""
	}

	email = strings.ToLower(strings.TrimSpace(email))

	for _, cp := range tracker.Checkpoints {

		for i, e := range cp.Emails {
			if strings.ToLower(strings.TrimSpace(e)) == email {

				// 🔐 pastikan address ada
				if i < len(cp.Addresses) {
					return cp.Addresses[i]
				}
			}
		}
	}

	return ""
}
