package services

import (
	"doc-tracker/mempool"
	"doc-tracker/utils"
	"fmt"
	"time"
)

func UpdateCheckpointStatus(trackerID string, checkpointAddr string, evidenceHash string, evidencePath string) error {
	tracker := mempool.GetByID(trackerID)
	if tracker == nil {
		return fmt.Errorf("tracker not found")
	}
	// fmt.Printf("Updating checkpoint %s for tracker %s with evidence hash %s and path %s\n", checkpointAddr, trackerID, evidenceHash, evidencePath)

	updated := false

	// fmt.Printf("==>Checking checkpoint %s for tracker %s\n", checkpointAddr, trackerID)
	// Update matching checkpoint
	for i, cp := range tracker.Checkpoints {
		// fmt.Printf("==>Checking checkpoint %s at index %d\n", cp.Address, i)
		if cp.Address == checkpointAddr {
			// fmt.Printf("<<<< Found matching checkpoint %s at index %d\n", cp.Address, i)
			if cp.IsCompleted {
				// fmt.Printf("]]] Checkpoint %s for tracker %s is already completed\n", cp.Address, trackerID)
				return fmt.Errorf("checkpoint %s for tracker %s is already completed", cp.Address, trackerID)
			}
			tracker.Checkpoints[i].IsCompleted = true
			tracker.Checkpoints[i].CompletedAt = time.Now().Unix()
			tracker.Checkpoints[i].Note = cp.Note

			publicKeyStr, err := utils.LoadKeys()
			if err != nil {
				// fmt.Printf("failed to get public key: %v\n", err)
				continue
			}
			// fmt.Printf("1) Public key loaded for checkpoint %s: %s\n", cp.Address, publicKeyStr.PublicKey)

			// Enkripsi data
			sharedSecret, err := utils.DeriveSharedSecret(publicKeyStr.PrivateKey, publicKeyStr.PublicKey)
			if err != nil {
				return fmt.Errorf("failed to derive secret: %v", err)
			}
			// fmt.Printf("2) Shared secret derived for checkpoint %s\n", cp.Address)

			encKey := utils.GenerateEncryptionKey(sharedSecret)
			encryptedNote, err := utils.EncryptData(encKey, []byte(cp.Note))
			if err != nil {
				return fmt.Errorf("encryption failed: %v", err)
			}
			// fmt.Printf("3) Note encrypted for checkpoint %s\n", cp.Address)
			// Simpan encrypted note dan evidence
			tracker.Checkpoints[i].EncryptedNote = string(encryptedNote)

			tracker.Checkpoints[i].EvidenceHash = evidenceHash
			tracker.Checkpoints[i].EvidencePath = evidencePath
			updated = true
			// fmt.Printf("4) Evidence hash and path updated for checkpoint %s\n", cp.Address)
			break
		}
		// fmt.Printf("Checkpoint %s for tracker %s is not matched\n", cp.Address, trackerID)
	}

	if !updated {
		return fmt.Errorf("checkpoint not found or already completed")
	}
	// fmt.Printf("Checkpoint %s for tracker %s updated with evidence hash %s and path %s\n", checkpointAddr, trackerID, evidenceHash, evidencePath)

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
	mempool.SaveToFile()

	// Jika complete, broadcast ke miner
	if allComplete {
		go BroadcastToMempool(tracker)
	}

	return nil
}

func GetCheckpointAddressByEmail(trackerID string, email string) string {
	tracker := mempool.GetByID(trackerID)
	if tracker == nil {
		// fmt.Printf("Tracker %s not found\n", trackerID)
		return ""
	}

	for _, cp := range tracker.Checkpoints {
		if cp.Email == email {
			return cp.Address
		}
	}
	// fmt.Printf("Checkpoint address for email %s not found in tracker %s\n", email, trackerID)
	return ""
}
