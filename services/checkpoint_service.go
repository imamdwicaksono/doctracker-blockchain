package services

import (
	"doc-tracker/storage"
	"doc-tracker/utils"
	"fmt"
	"time"
)

func UpdateCheckpointStatus(trackerID string, checkpointAddr string, evidenceHash string, evidencePath string) error {
	tracker := storage.FindTrackerByID(trackerID)
	if tracker == nil {
		return fmt.Errorf("tracker not found")
	}

	updated := false

	// Update matching checkpoint
	for i, cp := range tracker.Checkpoints {
		if cp.Address == checkpointAddr && !cp.IsCompleted {
			tracker.Checkpoints[i].IsCompleted = true
			tracker.Checkpoints[i].CompletedAt = time.Now().Unix()
			tracker.Checkpoints[i].Note = cp.Note

			wallet, _ := GetWalletByEmail(cp.Email)
			publicKeyStr, err := utils.PublicKeyToString(wallet.PublicKey)
			if err != nil {
				return fmt.Errorf("failed to convert public key: %v", err)
			}
			encryptedNote, _ := utils.EncryptNoteRSA(cp.Note, publicKeyStr)
			// Simpan encrypted note dan evidence
			tracker.Checkpoints[i].EncryptedNote = encryptedNote

			tracker.Checkpoints[i].EvidenceHash = evidenceHash
			tracker.Checkpoints[i].EvidencePath = evidencePath
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("checkpoint not found or already completed")
	}

	// Cek semua checkpoint complete
	allComplete := true
	for _, cp := range tracker.Checkpoints {
		if !cp.IsCompleted {
			allComplete = false
			break
		}
	}

	// Update tracker status
	if allComplete {
		tracker.Status = "complete"
	} else {
		tracker.Status = "progress"
	}

	// Update in storage
	storage.UpdateTracker(tracker)

	// Jika complete, broadcast ke miner
	if allComplete {
		go BroadcastToMempool(tracker)
	}

	return nil
}
