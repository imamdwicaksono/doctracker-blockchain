package services

import (
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/p2p"
	"doc-tracker/utils"
	"time"

	"github.com/google/uuid"
)

var trackerMap = map[string]models.Tracker{}

const trackerFile = "data/trackers.json"

func init() {
	LoadTrackersFromDisk()
}

func CreateTracker(input models.Tracker) (models.Tracker, error) {
	input.ID = uuid.New().String()
	input.CreatedAt = time.Now().Unix()
	input.Status = "progress"

	// Generate wallet/address untuk pengaju
	senderWallet := GetOrCreateWallet(input.Creator)
	input.CreatorAddr = senderWallet.Address

	// Enkripsi checkpoint jika perlu
	for i, cp := range input.Checkpoints {
		receiverWallet := GetOrCreateWallet(cp.Email)
		input.Checkpoints[i].Address = receiverWallet.Address

		// Jika boleh melihat isi dokumen, enkripsi
		if cp.IsViewable {
			encrypted := utils.EncryptWithPublicKey(cp.Note, receiverWallet.PublicKey)
			input.Checkpoints[i].EncryptedNote = encrypted
			input.Checkpoints[i].Note = "" // kosongkan untuk keamanan
		} else {
			// Tidak boleh melihat, kosongkan isinya
			input.Checkpoints[i].Note = ""
			input.Checkpoints[i].EncryptedNote = ""
		}
	}

	// Simpan ke mempool dan broadcast
	mempool.Add(&input)
	p2p.BroadcastToMempool(input)

	return input, nil
}

func GetTrackers() ([]models.Tracker, error) {
	var trackers []models.Tracker
	err := mempool.Iterate(func(tx *models.Tracker) error {
		trackers = append(trackers, *tx)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return trackers, nil
}
func GetTrackerByID(id string) (models.Tracker, error) {
	var tracker models.Tracker
	err := mempool.Iterate(func(tx *models.Tracker) error {
		if tx.ID == id {
			tracker = *tx
			return nil // stop iterating
		}
		return nil
	})
	if err != nil {
		return models.Tracker{}, err
	}
	if tracker.ID == "" {
		return models.Tracker{}, utils.ErrNotFound
	}
	return tracker, nil
}
func GetTrackersByAddress(address string) ([]models.Tracker, error) {
	var trackers []models.Tracker
	err := mempool.Iterate(func(tx *models.Tracker) error {
		if tx.CreatorAddr == address || hasCheckpointForAddress(tx, address) {
			trackers = append(trackers, *tx)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return trackers, nil
}

func hasCheckpointForAddress(tx *models.Tracker, address string) bool {
	for _, cp := range tx.Checkpoints {
		if cp.Address == address {
			return true
		}
	}
	return false
}

func SaveTrackersToDisk() {
	utils.SaveToFile(trackerFile, trackerMap)
}

func LoadTrackersFromDisk() {
	utils.LoadFromFile(trackerFile, &trackerMap)
}
